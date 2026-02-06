package models

type Company struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Description      *string `json:"description,omitempty"`
	ParentCompanyID  *int    `json:"parent_company_id,omitempty"`
	Industry         *string `json:"industry,omitempty"`
	SubIndustry      *string `json:"sub_industry,omitempty"`
	Headquarters     *string `json:"headquarters,omitempty"`
	DateIncorporated *string `json:"date_incorporated,omitempty"`
	UserID           *int    `json:"user_id,omitempty"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

type Post struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	CompanyID   *int    `json:"company_id,omitempty"`
	UserID      int     `json:"user_id"`
	Upvotes     int     `json:"upvotes"`
	Downvotes   int     `json:"downvotes"`
}

type Comment struct {
	ID        int    `json:"id"`
	Message   string `json:"message"`
	PostID    int    `json:"post_id"`
	UserID    int    `json:"user_id"`
	Upvotes   int    `json:"upvotes"`
	Downvotes int    `json:"downvotes"`
}
