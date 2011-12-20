package librato

import (
	"fmt"
	"json"
	"os"
)

type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	APIToken  string `json:"api_token"`
	Reference string `json:"reference"`
	Name      string `json:"name"`
	Company   string `json:"company"`
	Country   string `json:"country"`
	TimeZone  string `json:"time_zone"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UsersResponse struct {
	Query QueryResponse `json:"query"`
	Users []User        `json:"users"`
}

func (u *User) String() string {
	return fmt.Sprintf("{ID:%d Email:%s}", u.ID, u.Email)
}

func (r *UsersResponse) String() string {
	return fmt.Sprintf("{Query:%s Users:%s}", r.Query.String(), r.Users)
}

func (met *Metrics) GetUsers(reference string, email string) (*UsersResponse, os.Error) {
	res, err := met.get(librato_metrics_users_api_url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, os.NewError(res.Status)
	}

	var users UsersResponse
	jdec := json.NewDecoder(res.Body)
	err = jdec.Decode(&users)
	return &users, err
}
