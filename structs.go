package main

type MyTransport struct{}

type WriteCounter struct {
	Total      uint64
	TotalStr   string
	Downloaded uint64
	Percentage int
}

type Config struct {
	Email            string
	Password         string
	Urls             []string
	Format           int
	OutPath          string
	TrackTemplate    string
	DownloadBooklets bool
	MaxCoverSize     bool
	KeepCover        bool
}

type Args struct {
	Urls    []string `arg:"positional, required"`
	Format  int      `arg:"-f" default:"-1"`
	OutPath string   `arg:"-o"`
}

type Auth struct {
	AccessToken  string `json:"accessToken"`
	SessionToken string `json:"sessionToken"`
	User         struct {
		ID              string      `json:"id"`
		OldID           int         `json:"old_id"`
		Email           string      `json:"email"`
		Firstname       interface{} `json:"firstname"`
		EncryptedEmail  string      `json:"encrypted_email"`
		Premium         bool        `json:"premium"`
		Plan            string      `json:"plan"`
		PlanDisplayName string      `json:"plan_display_name"`
		PlanBeginsAt    int         `json:"plan_begins_at"`
		PlanEndsAt      int         `json:"plan_ends_at"`
		PlanRenewsAt    int         `json:"plan_renews_at"`
		PlanCanceledAt  interface{} `json:"plan_canceled_at"`
		PreviousPlan    interface{} `json:"previous_plan"`
		FromID          string      `json:"from_id"`
		Flow            interface{} `json:"flow"`
		Preferences     struct {
			WebAudioQuality     int    `json:"web_audio_quality"`
			WebTheme            string `json:"web_theme"`
			IosAudioQuality     int    `json:"ios_audio_quality"`
			AndroidAudioQuality int    `json:"android_audio_quality"`
			Locale              string `json:"locale"`
			MarketingOptIn      bool   `json:"marketing_opt_in"`
			ResearchOptIn       bool   `json:"research_opt_in"`
		} `json:"preferences"`
		PlanSubscriptionType   interface{} `json:"plan_subscription_type"`
		PlanProductID          interface{} `json:"plan_product_id"`
		HadAppleSubscription   bool        `json:"had_apple_subscription"`
		HadAndroidSubscription bool        `json:"had_android_subscription"`
		HadStripeSubscription  bool        `json:"had_stripe_subscription"`
		HadTrialOptIn          bool        `json:"had_trial_opt_in"`
		HadTrialOptOut         bool        `json:"had_trial_opt_out"`
		Subscription           struct {
		} `json:"subscription"`
		HadTrial          bool          `json:"had_trial"`
		PostWebEOFUser    bool          `json:"post_web_eof_user"`
		ShouldOfferOptOut bool          `json:"should_offer_opt_out"`
		EffectiveCountry  string        `json:"effective_country"`
		LastIPCountry     string        `json:"last_ip_country"`
		Tier              interface{}   `json:"tier"`
		PaymentCountry    interface{}   `json:"payment_country"`
		Credits           []interface{} `json:"credits"`
		CreatedAt         int           `json:"created_at"`
		IDSignature       string        `json:"id_signature"`
		LoginMethods      []string      `json:"login_methods"`
		StripeDiscount    interface{}   `json:"stripe_discount"`
		StripeStatus      interface{}   `json:"stripe_status"`
		Tags              []interface{} `json:"tags"`
		Features          struct {
			FullCatalogue                 bool  `json:"full_catalogue"`
			StreamingToAllExternalPlayers bool  `json:"streaming_to_all_external_players"`
			AudioQuality                  []int `json:"audio_quality"`
			LosslessAsGift                bool  `json:"lossless_as_gift"`
			AllowPlaybackSkip             bool  `json:"allow_playback_skip"`
			SkipLimits                    struct {
				Enabled              bool  `json:"enabled"`
				MaxSkips             int   `json:"max_skips"`
				ResetIntervalSeconds int   `json:"reset_interval_seconds"`
				TooltipPositions     []int `json:"tooltip_positions"`
			} `json:"skip_limits"`
			ChangePlaybackPosition bool `json:"change_playback_position"`
			PlayIndividualTracks   bool `json:"play_individual_tracks"`
			MaxSkips               int  `json:"max_skips"`
			Collections            struct {
				AccessLimit          int  `json:"access_limit"`
				ShowAlbumsPlayAll    bool `json:"show_albums_play_all"`
				ShowPlaylistsPlayAll bool `json:"show_playlists_play_all"`
			} `json:"collections"`
			Gch struct {
				AllowConcertPlayback bool `json:"allow_concert_playback"`
			} `json:"gch"`
			IntervalBetweenIntermission int  `json:"interval_between_intermission"`
			OfflineListening            bool `json:"offline_listening"`
			UserCollections             bool `json:"user_collections"`
			UserPlaylists               bool `json:"user_playlists"`
		} `json:"features"`
		SheerIDVerified interface{} `json:"sheer_id_verified"`
		Tickets         []struct {
			TicketType  string `json:"ticket_type"`
			TicketID    string `json:"ticket_id"`
			ProductType string `json:"product_type"`
			ProductID   string `json:"product_id"`
			ExpiresAt   int    `json:"expires_at"`
		} `json:"tickets"`
	} `json:"user"`
}

type AlbumMeta struct {
	Result struct {
		LongDescription  string         `json:"longDescription"`
		Copyright        string         `json:"copyright"`
		PublishDate      string         `json:"publishDate"`
		InterestGroupIds []string       `json:"interestGroupIds"`
		UPC              string         `json:"upc"`
		Description      string         `json:"description"`
		BookletURL       string         `json:"bookletUrl"`
		Title            string         `json:"title"`
		UUID             string         `json:"uuid"`
		Tags             []string       `json:"tags"`
		TrackIds         []string       `json:"trackIds"`
		CopyrightYear    int            `json:"copyrightYear"`
		Slug             string         `json:"slug"`
		Participants     []Participants `json:"participants`
		ID               string         `json:"id"`
		ImageURL         string         `json:"imageUrl"`
		Tracks           []TrackMeta    `json:"tracks"`
	} `json:"result"`
}

type Participants struct {
	Participation  float64 `json:"participation"`
	NumOccurrences int     `json:"numOccurrences"`
	Person         struct {
		Forename      string `json:"forename"`
		InstrumentKey string `json:"_instrumentKey"`
		InstrumentID  int    `json:"instrumentId"`
		Surname       string `json:"surname"`
		Instrument    string `json:"instrument"`
	} `json:"person"`
	Popularity float64     `json:"popularity"`
	Name       string      `json:"name"`
	Ensemble   interface{} `json:"ensemble"`
	ID         int         `json:"id"`
	Type       string      `json:"type"`
	Key        string      `json:"_key"`
}

type Authors struct {
	Persons []struct {
		Forename string `json:"forename"`
		Surname  string `json:"surname"`
		Name     string `json:"name"`
		ID       int    `json:"id"`
	} `json:"persons"`
	AuthorType string `json:"authorType"`
	Key        string `json:"key"`
}

type TrackMeta struct {
	Duration int `json:"duration"`
	Piece    struct {
		PopularTitle       interface{} `json:"popularTitle"`
		Number             interface{} `json:"number"`
		StructuralSubtitle interface{} `json:"structuralSubtitle"`
		Workpart           struct {
			Work struct {
				PopularTitle interface{} `json:"popularTitle"`
				Composer     struct {
					PlayEventCount   int           `json:"playEventCount"`
					AlternativeNames []interface{} `json:"alternativeNames"`
					Forename         string        `json:"forename"`
					Kind             string        `json:"kind"`
					Surname          string        `json:"surname"`
					Popularity       float64       `json:"popularity"`
					Name             string        `json:"name"`
					ID               int           `json:"id"`
					Key              string        `json:"_key"`
					Slugs            interface{}   `json:"slugs"`
				} `json:"composer"`
				Epoch struct {
					ID    int    `json:"id"`
					Title string `json:"title"`
					Key   string `json:"_key"`
				} `json:"epoch"`
				Tonality struct {
				} `json:"tonality"`
				Key                string      `json:"_key"`
				Title              string      `json:"title"`
				Number             string      `json:"number"`
				StructuralSubtitle interface{} `json:"structuralSubtitle"`
				DefaultRecordingID int         `json:"defaultRecordingId"`
				OpusNumber         string      `json:"opusNumber"`
				Subgenre           struct {
					ID    int    `json:"id"`
					Title string `json:"title"`
					Key   string `json:"_key"`
				} `json:"subgenre"`
				Subtitle interface{} `json:"subtitle"`
				Genre    struct {
					ID    int    `json:"id"`
					Title string `json:"title"`
					Key   string `json:"_key"`
				} `json:"genre"`
				ID              int `json:"id"`
				Authors         []Authors
				CompositionYear string `json:"compositionYear"`
			} `json:"work"`
			ID         int         `json:"id"`
			Position   int         `json:"position"`
			Title      interface{} `json:"title"`
			IsOverture interface{} `json:"isOverture"`
		} `json:"workpart"`
		Subtitle        interface{} `json:"subtitle"`
		DescriptionName string      `json:"descriptionName"`
		ID              int         `json:"id"`
		Position        int         `json:"position"`
		Title           string      `json:"title"`
		PieceName       interface{} `json:"pieceName"`
	} `json:"piece"`
	Recording struct {
		Summary string      `json:"summary"`
		Venue   interface{} `json:"venue"`
		Albums  []string    `json:"albums"`
		Work    struct {
			ID int `json:"id"`
		} `json:"work"`
		CreatedAt     int64 `json:"created_at"`
		RecordingDate struct {
			FromKeyword interface{} `json:"fromKeyword"`
			From        string      `json:"from"`
			To          interface{} `json:"to"`
		} `json:"recordingDate"`
		Conductor struct {
			PlayEventCount   int           `json:"playEventCount"`
			AlternativeNames []interface{} `json:"alternativeNames"`
			Forename         string        `json:"forename"`
			Kind             string        `json:"kind"`
			Surname          string        `json:"surname"`
			Popularity       float64       `json:"popularity"`
			Name             string        `json:"name"`
			ID               int           `json:"id"`
			Key              string        `json:"_key"`
			Slugs            interface{}   `json:"slugs"`
		} `json:"conductor"`
		InterestGroups []int       `json:"interestGroups"`
		IsExclusive    bool        `json:"isExclusive"`
		RecordingType  interface{} `json:"recordingType"`
		Location       struct {
			Name string `json:"name"`
			ID   int    `json:"id"`
		} `json:"location"`
		ID        int           `json:"id"`
		IsPremium bool          `json:"isPremium"`
		Soloists  []interface{} `json:"soloists"`
		Ensembles []struct {
			PlayEventCount   int           `json:"playEventCount"`
			AlternativeNames []interface{} `json:"alternativeNames"`
			Kind             string        `json:"kind"`
			Popularity       float64       `json:"popularity"`
			Name             string        `json:"name"`
			ID               int           `json:"id"`
			Key              string        `json:"_key"`
			Slugs            interface{}   `json:"slugs"`
		} `json:"ensembles"`
		Tags []interface{} `json:"tags"`
	} `json:"recording"`
	ID       int `json:"id"`
	Position int `json:"position"`
}

type FileMeta struct {
	Length int
	URL    string
}
