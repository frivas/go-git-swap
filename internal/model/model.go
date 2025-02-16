package model

type GitProfile struct {
    Name       string `json:"name"`       // Profile name
    Email      string `json:"email"`      // Git email
    Username   string `json:"username"`   // Git username
    FullName   string `json:"full_name"`  // Git user.name
    SigningKey string `json:"signing_key"`
}

type Profiles struct {
    Active   string                `json:"active"`
    Profiles map[string]GitProfile `json:"profiles"`
}
