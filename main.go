package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/bogem/id3v2"
	"github.com/dustin/go-humanize"
	"github.com/go-flac/flacpicture"
	"github.com/go-flac/flacvorbis"
	"github.com/go-flac/go-flac"
)

var (
	jar, _ = cookiejar.New(nil)
	client = &http.Client{Jar: jar, Transport: &MyTransport{}}
)

func (t *MyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add(
		"User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 "+
			"(KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36",
	)
	req.Header.Add(
		"Referer", "https://app.idagio.com/",
	)
	return http.DefaultTransport.RoundTrip(req)
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Downloaded += uint64(n)
	percentage := float64(wc.Downloaded) / float64(wc.Total) * float64(100)
	wc.Percentage = int(percentage)
	fmt.Printf("\r%d%%, %s/%s ", wc.Percentage, humanize.Bytes(wc.Downloaded), wc.TotalStr)
	return n, nil
}

func getScriptDir() (string, error) {
	var (
		ok    bool
		err   error
		fname string
	)
	if filepath.IsAbs(os.Args[0]) {
		_, fname, _, ok = runtime.Caller(0)
		if !ok {
			return "", errors.New("Failed to get script filename.")
		}
	} else {
		fname, err = os.Executable()
		if err != nil {
			return "", err
		}
	}
	scriptDir := filepath.Dir(fname)
	return scriptDir, nil
}

func readTxtFile(path string) ([]string, error) {
	var lines []string
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	return lines, nil
}

func contains(lines []string, value string) bool {
	for _, line := range lines {
		if strings.EqualFold(line, value) {
			return true
		}
	}
	return false
}

func processUrls(urls []string) ([]string, error) {
	var processed []string
	var txtPaths []string
	for _, url := range urls {
		if strings.HasSuffix(url, ".txt") && !contains(txtPaths, url) {
			txtLines, err := readTxtFile(url)
			if err != nil {
				return nil, err
			}
			for _, txtLine := range txtLines {
				if !contains(processed, txtLine) {
					processed = append(processed, txtLine)
				}
			}
			txtPaths = append(txtPaths, url)
		} else {
			if !contains(processed, url) {
				processed = append(processed, url)
			}
		}
	}
	return processed, nil
}

func readConfig() (*Config, error) {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		return nil, err
	}
	var obj Config
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

func parseArgs() *Args {
	var args Args
	arg.MustParse(&args)
	return &args
}

func parseCfg() (*Config, error) {
	resolveFormat := map[int]int{
		1: 50,
		2: 70,
		3: 90,
	}
	cfg, err := readConfig()
	if err != nil {
		return nil, err
	}
	args := parseArgs()
	if args.Format != -1 {
		cfg.Format = args.Format
	}
	if !(cfg.Format >= 1 && cfg.Format <= 3) {
		return nil, errors.New("Format must be between 1 and 3.")
	}
	cfg.Format = resolveFormat[cfg.Format]
	if args.OutPath != "" {
		cfg.OutPath = args.OutPath
	}
	if cfg.OutPath == "" {
		cfg.OutPath = "Idagio downloads"
		// Use $HOME/Music if available
		homeDir, _ := os.UserHomeDir()
		musicPath := path.Join(homeDir, "Music")
		if available, _ := isDirectory(musicPath); available {
			cfg.OutPath = musicPath
		}
		fmt.Printf("Download path set: %s\n", cfg.OutPath)
	}
	cfg.Urls, err = processUrls(args.Urls)
	if err != nil {
		errString := fmt.Sprintf("Failed to process URLs.\n %s", err)
		return nil, errors.New(errString)
	}
	return cfg, nil
}

func auth(email, pwd string) (string, string, error) {
	url := "https://app.idagio.com/login.json"
	postMap := map[string]string{
		"username": email,
		"password": pwd,
		"type":     "Modal",
	}
	postData, err := json.Marshal(postMap)
	if err != nil {
		return "", "", err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(postData))
	if err != nil {
		return "", "", err
	}
	req.Header.Add("Content-Type", "application/json")
	do, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer do.Body.Close()
	if do.StatusCode != http.StatusOK {
		return "", "", errors.New(do.Status)
	}
	var obj Auth
	if err != nil {
		return "", "", err
	}
	err = json.NewDecoder(do.Body).Decode(&obj)
	if err != nil {
		return "", "", err
	}
	if !strings.HasPrefix(obj.User.Plan, "Premium") {
		return "", "", errors.New("User doesn't have a Premium plan.")
	}
	return obj.AccessToken, obj.User.PlanDisplayName, nil
}

func resolve(format int) (string, string) {
	resolveMap := map[int][2]string{
		50: [2]string{".mp3", "192 Kbps MP3"},
		70: [2]string{".mp3", "320 Kbps MP3"},
		90: [2]string{".flac", "16-bit / 44.1 kHz FLAC"},
	}
	resolved := resolveMap[format]
	return resolved[0], resolved[1]
}

func getMeta(albumId, token string) (*AlbumMeta, error) {
	url := "https://api.idagio.com/v2.0/albums/" + albumId
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "bearer "+token)
	do, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer do.Body.Close()
	if do.StatusCode != http.StatusOK {
		return nil, errors.New(do.Status)
	}
	var obj AlbumMeta
	err = json.NewDecoder(do.Body).Decode(&obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

func checkUrl(url string) string {
	url = strings.Split(url, "?")[0]
	regexString := `^https://(?:app|www).idagio.com/albums/([a-zA-Z\d-]+)$`
	regex := regexp.MustCompile(regexString)
	match := regex.FindStringSubmatch(url)
	if match == nil {
		return ""
	}
	return match[1]
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func makeDir(path string) error {
	err := os.Mkdir(path, 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}
	return nil
}

func fileExists(path string) (bool, error) {
	f, err := os.Stat(path)
	if err == nil {
		return !f.IsDir(), nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func multiParticipants(meta []Participants) (string, string) {
	var (
		composer     string
		participants string
	)
	for _, participant := range meta {
		if participant.Type == "composer" {
			composer = participant.Name
		}
		participants += participant.Name + ", "
	}
	return participants[:len(participants)-2], composer
}

func multiAuthors(meta []Authors) string {
	var authors string
	for _, author := range meta {
		authors += author.Persons[0].Name + ", "
	}
	return authors[:len(authors)-2]
}

func getComposer(meta []Participants) string {
	for _, author := range meta {
		if author.Type == "composer" {
			return author.Name
		}
	}
	return ""
}

func parseAlbumMeta(meta *AlbumMeta) map[string]string {
	albumArtists, composer := multiParticipants(meta.Result.Participants)
	parsedMeta := map[string]string{
		"album":       meta.Result.Title,
		"albumArtist": albumArtists,
		"composer":    composer,
		"copyright":   meta.Result.Copyright,
		"upc":         meta.Result.UPC,
		"year":        strings.Split(meta.Result.PublishDate, "-")[0],
	}
	return parsedMeta
}

func parseTrackMeta(meta *TrackMeta, albMeta map[string]string, trackNum, trackTotal int) map[string]string {
	artists := multiAuthors(meta.Piece.Workpart.Work.Authors)
	albMeta["artist"] = artists
	albMeta["genre"] = meta.Piece.Workpart.Work.Genre.Title
	albMeta["title"] = meta.Piece.Title
	albMeta["track"] = strconv.Itoa(trackNum)
	albMeta["trackPad"] = fmt.Sprintf("%02d", trackNum)
	albMeta["trackTotal"] = strconv.Itoa(trackTotal)
	return albMeta
}

func sanitize(filename string) string {
	regex := regexp.MustCompile(`[\/:*?"><|]`)
	sanitized := regex.ReplaceAllString(filename, "_")
	return sanitized
}

func downloadCover(upc, path string, maxCover bool) error {
	url := "https://idagio-images.global.ssl.fastly.net/albums/" + upc + "/main.jpg"
	if !maxCover {
		url += "?h=600"
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	req, err := client.Get(url)
	if err != nil {
		return err
	}
	defer req.Body.Close()
	if req.StatusCode != http.StatusOK {
		return errors.New(req.Status)
	}
	_, err = io.Copy(f, req.Body)
	return err
}

func parseTemplate(templateText string, tags map[string]string) string {
	var buffer bytes.Buffer
	for {
		err := template.Must(template.New("").Parse(templateText)).Execute(&buffer, tags)
		if err == nil {
			break
		}
		fmt.Println("Failed to parse template. Default will be used instead.")
		templateText = "{{.trackPad}}. {{.title}}"
		buffer.Reset()
	}
	return buffer.String()
}

func getStreamUrl(token string, trackId, format int) (string, error) {
	_url := "https://api.idagio.com/v1.8/content/track/" + strconv.Itoa(trackId)
	req, err := http.NewRequest(http.MethodGet, _url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "bearer "+token)
	query := url.Values{}
	query.Set("quality", strconv.Itoa(format))
	query.Set("format", "2")
	query.Set("client_type", "web-3")
	query.Set("client_version", "0.0.123")
	query.Set("device_id", "web")
	req.URL.RawQuery = query.Encode()
	do, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer do.Body.Close()
	if do.StatusCode != http.StatusOK {
		return "", errors.New(do.Status)
	}
	var obj FileMeta
	err = json.NewDecoder(do.Body).Decode(&obj)
	if err != nil {
		return "", err
	}
	return obj.URL, nil
}

func downloadTrack(trackPath, url string) (string, string, error) {
	f, err := os.OpenFile(trackPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return "", "", err
	}
	defer f.Close()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Add("Range", "bytes=0-")
	do, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer do.Body.Close()
	if do.StatusCode != http.StatusOK && do.StatusCode != http.StatusPartialContent {
		return "", "", errors.New(do.Status)
	}
	totalBytes := uint64(do.ContentLength)
	counter := &WriteCounter{Total: totalBytes, TotalStr: humanize.Bytes(totalBytes)}
	_, err = io.Copy(f, io.TeeReader(do.Body, counter))
	fmt.Println("")
	baseKeyIv := do.Header["X-X"][0]
	baseKeyIvSplit := strings.Split(baseKeyIv, " ")
	return baseKeyIvSplit[0], baseKeyIvSplit[1], err
}

// Not sure if always encrypted.
func decryptTrack(baseKey, iv, inPath, outPath string) error {
	var stream cipher.Stream
	fmt.Println("Decrypting...")
	in_f, err := os.Open(inPath)
	if err != nil {
		return err
	}
	// No defer close so it can be deleted later.
	key := []byte(baseKey + "prod-media-c-YaiJaoni7iebeed5")
	rawKey := sha256.Sum256(key)
	hexKey := hex.EncodeToString(rawKey[:8])
	block, err := aes.NewCipher([]byte(hexKey))
	if err != nil {
		in_f.Close()
		return err
	}
	stream = cipher.NewCTR(block, []byte(iv))
	reader := &cipher.StreamReader{S: stream, R: in_f}
	out_f, err := os.Create(outPath)
	if err != nil {
		in_f.Close()
		return err
	}
	defer out_f.Close()
	_, err = io.Copy(out_f, reader)
	if err != nil {
		in_f.Close()
		return err
	}
	in_f.Close()
	err = os.Remove(inPath)
	if err != nil {
		fmt.Println("Failed to delete encrypted track.")
	}
	return nil
}

func writeFlacTags(decTrackPath string, tags map[string]string, imgData []byte) error {
	f, err := flac.ParseFile(decTrackPath)
	if err != nil {
		return err
	}
	tag := flacvorbis.New()
	for k, v := range tags {
		tag.Add(strings.ToUpper(k), v)
	}
	tagMeta := tag.Marshal()
	f.Meta = append(f.Meta, &tagMeta)
	if imgData != nil {
		picture, err := flacpicture.NewFromImageData(
			flacpicture.PictureTypeFrontCover, "", imgData, "image/jpeg",
		)
		if err != nil {
			return err
		}
		pictureMeta := picture.Marshal()
		f.Meta = append(f.Meta, &pictureMeta)
	}
	err = f.Save(decTrackPath)
	return err
}

func writeMp3Tags(decTrackPath string, tags map[string]string, imgData []byte) error {
	tags["track"] += "/" + tags["trackTotal"]
	resolve := map[string]string{
		"album":       "TALB",
		"artist":      "TPE1",
		"albumArtist": "TPE2",
		"composer":    "TCOM",
		"copyright":   "TCOP",
		"genre":       "TCON",
		"title":       "TIT2",
		"track":       "TRCK",
		"year":        "TYER",
	}
	tag, err := id3v2.Open(decTrackPath, id3v2.Options{Parse: true})
	if err != nil {
		return err
	}
	defer tag.Close()
	for k, v := range tags {
		resolved, ok := resolve[k]
		if ok {
			tag.AddTextFrame(resolved, tag.DefaultEncoding(), v)
		}
	}
	if imgData != nil {
		imgFrame := id3v2.PictureFrame{
			Encoding:    id3v2.EncodingUTF8,
			MimeType:    "image/jpeg",
			PictureType: id3v2.PTFrontCover,
			Picture:     imgData,
		}
		tag.AddAttachedPicture(imgFrame)
	}
	err = tag.Save()
	return err
}

func writeTags(decTrackPath, coverPath, ext string, tags map[string]string) error {
	var (
		err     error
		imgData []byte
	)
	if coverPath != "" {
		imgData, err = ioutil.ReadFile(coverPath)
		if err != nil {
			return err
		}
	}
	delete(tags, "trackPad")
	if ext == ".flac" {
		err = writeFlacTags(decTrackPath, tags, imgData)
	} else {
		err = writeMp3Tags(decTrackPath, tags, imgData)
	}
	return err
}

func downloadBooklet(path, url string) error {
	fmt.Println("Downloading booklet...")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	req, err := client.Get(url)
	if err != nil {
		return err
	}
	defer req.Body.Close()
	if req.StatusCode != http.StatusOK {
		return errors.New(req.Status)
	}
	_, err = io.Copy(f, req.Body)
	return err
}

func init() {
	fmt.Println(`                                                                    
 _____ ____  _____ _____ _____ _____    ____                _           _         
|     |    \|  _  |   __|     |     |  |    \ ___ _ _ _ ___| |___ ___ _| |___ ___ 
|-   -|  |  |     |  |  |-   -|  |  |  |  |  | . | | | |   | | . | .'| . | -_|  _|
|_____|____/|__|__|_____|_____|_____|  |____/|___|_____|_|_|_|___|__,|___|___|_|
`)
	scriptDir, err := getScriptDir()
	if err != nil {
		panic(err)
	}
	err = os.Chdir(scriptDir)
	if err != nil {
		panic(err)
	}
}

func main() {
	cfg, err := parseCfg()
	if err != nil {
		errString := fmt.Sprintf("Failed to parse config file. %s", err)
		panic(errString)
	}
	token, plan, err := auth(cfg.Email, cfg.Password)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Signed in successfully - %s plan.\n\n", plan)
	err = makeDir(cfg.OutPath)
	if err != nil {
		errString := fmt.Sprintf("Failed to make output folder. %s", err)
		panic(errString)
	}
	ext, specs := resolve(cfg.Format)
	albumTotal := len(cfg.Urls)
	for albumNum, url := range cfg.Urls {
		fmt.Printf("Album %d of %d:\n", albumNum+1, albumTotal)
		albumSlug := checkUrl(url)
		if albumSlug == "" {
			fmt.Println("Invalid URL:", url)
			continue
		}
		meta, err := getMeta(albumSlug, token)
		if err != nil {
			fmt.Println("Failed to get metadata.\n", err)
			continue
		}
		parsedAlbMeta := parseAlbumMeta(meta)
		albFolder := parsedAlbMeta["composer"] + " - " + parsedAlbMeta["album"]
		fmt.Println(albFolder)
		if len(albFolder) > 120 {
			fmt.Println("Album folder was chopped as it exceeds 120 characters.")
			albFolder = albFolder[:120]
		}
		albumPath := filepath.Join(cfg.OutPath, sanitize(albFolder))
		err = makeDir(albumPath)
		if err != nil {
			fmt.Println("Failed to make album folder.\n", err)
			continue
		}
		coverPath := filepath.Join(albumPath, "cover.jpg")
		err = downloadCover(parsedAlbMeta["upc"], coverPath, cfg.MaxCoverSize)
		if err != nil {
			fmt.Println("Failed to get cover.\n", err)
			coverPath = ""
		}
		bookletUrl := meta.Result.BookletURL
		if bookletUrl != "" && cfg.DownloadBooklets {
			bookletPath := filepath.Join(albumPath, "booklet.pdf")
			err = downloadBooklet(bookletPath, bookletUrl)
			if err != nil {
				fmt.Println("Failed to download booklet.\n", err)
			}
		}
		trackTotal := len(meta.Result.Tracks)
		for trackNum, track := range meta.Result.Tracks {
			trackNum++
			parsedMeta := parseTrackMeta(&track, parsedAlbMeta, trackNum, trackTotal)
			trackFname := parseTemplate(cfg.TrackTemplate, parsedMeta)
			trackPath := filepath.Join(albumPath, sanitize(trackFname)+ext)
			exists, err := fileExists(trackPath)
			if err != nil {
				fmt.Println("Failed to check if track already exists locally.\n", err)
				continue
			}
			if exists {
				fmt.Println("Track already exists locally.")
				continue
			}
			streamUrl, err := getStreamUrl(token, track.ID, cfg.Format)
			if err != nil {
				fmt.Println("Failed to get track stream URL.\n", err)
				continue
			}
			encTrackPath := filepath.Join(albumPath, strconv.Itoa(trackNum)+".encrypted")
			decTrackPath := filepath.Join(albumPath, strconv.Itoa(trackNum)+ext)
			fmt.Printf("Downloading track %d of %d: %s - %s\n", trackNum, trackTotal, parsedMeta["title"], specs)
			baseKey, iv, err := downloadTrack(encTrackPath, streamUrl)
			if err != nil {
				fmt.Println("Failed to download track.\n", err)
				continue
			}
			err = decryptTrack(baseKey, iv, encTrackPath, decTrackPath)
			if err != nil {
				fmt.Println("Failed to decrypt track.\n", err)
			}
			err = writeTags(decTrackPath, coverPath, ext, parsedMeta)
			if err != nil {
				fmt.Println("Failed to write tags.\n", err)
				continue
			}
			err = os.Rename(decTrackPath, trackPath)
			if err != nil {
				fmt.Println("Failed to rename track. Left as is.")
			}
		}
		if coverPath != "" && !cfg.KeepCover {
			err = os.Remove(coverPath)
			if err != nil {
				fmt.Println("Failed to delete cover.\n", err)
			}
		}
	}
}
