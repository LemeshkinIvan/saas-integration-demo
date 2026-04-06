package oauth_models

type AmoAccountOutput struct {
	AmoID     string
	Name      string
	Subdomain string

	AmojoID string
}

// by account subdomain == referer
type AmoAccount struct {
	AmoID     int    `json:"id"`
	Name      string `json:"name"`
	Subdomain string `json:"subdomain"`

	AmojoID string `json:"amojo_id"`

	// Language string `json:"language"`
	// Country  string `json:"country"`

	// Currency       string `json:"currency"`
	// CurrencySymbol string `json:"currency_symbol"`

	// IsHelpbotEnabled   bool `json:"is_helpbot_enabled"`
	// IsTechnicalAccount bool `json:"is_technical_account"`
}

/*
	FULL AMO RESPONSE
{
    "id": 32499146,
    "name": "daos.tech.widgets",
    "subdomain": "daostechwidgets",
    "language": "ru",
    "created_at": 1759077365,
    "created_by": 0,
    "updated_at": 1759077365,
    "updated_by": 0,
    "current_user_id": 12649650,
    "country": "RS",
    "currency": "RSD",
    "currency_symbol": "дин",
    "customers_mode": "disabled",
    "is_unsorted_on": true,
    "mobile_feature_version": 0,
    "is_loss_reason_enabled": true,
    "is_helpbot_enabled": false,
    "is_technical_account": true,
    "contact_name_display_order": 1,
    "amojo_id": "c906ccae-c143-46e2-a164-76b4c95eb1a9",
    "_links": {
        "self": {
            "href": "https://daostechwidgets.amocrm.ru/api/v4/account"
        }
    }
}
*/
