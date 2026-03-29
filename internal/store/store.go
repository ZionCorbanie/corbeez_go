package store

import "time"

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	UserType    string    `json:"user_type" gorm:"type:enum('admin','user');default:user"`
}

type Session struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	SessionID string `json:"session_id"`
	UserID    uint   `json:"user_id"`
	User      User   `gorm:"foreignKey:UserID" json:"user"`
}

type Poll struct {
    ID      uint         `gorm:"primaryKey" json:"id"`
    Title   string       `json:"title" gorm:"type:varchar(255)"`
    Options []PollOption `gorm:"foreignKey:PollID;constraint:OnDelete:CASCADE;" json:"options"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

type PollOption struct {
    ID     uint       `gorm:"primaryKey" json:"id"`
    PollID uint       `json:"poll_id"`
    Poll   Poll       `gorm:"foreignKey:PollID" json:"poll"`
    Option string     `json:"option" gorm:"type:varchar(255)"`
}

type PollVote struct {
    ID       uint      `gorm:"primaryKey" json:"id"`
    UserID   uint      `json:"user_id"`
    User     User      `gorm:"foreignKey:UserID" json:"user"`
    OptionID uint      `json:"option_id"`
    Option   PollOption `gorm:"foreignKey:OptionID" json:"option"`
    PollID uint       `json:"poll_id"`
    Poll   Poll       `gorm:"foreignKey:PollID" json:"poll"`
	Time   time.Time `json:"time"`
}

type UserStore interface {
	CreateUser(email string, password string) error
	GetUser(email string) (*User, error)
}

type SessionStore interface {
	CreateSession(session *Session) (*Session, error)
	GetUserFromSession(sessionID string, userID string) (*User, error)
}

type PollStore interface {
    CreatePoll(poll *Poll) error
    GetPoll(pollId string) (*Poll, error)
    GetPolls() (*[]Poll, error)
    DeletePoll(pollId string) error
    PutPoll(poll Poll) error
    VotePoll(pollId uint, optionId uint, userId uint) error
    GetPollVotes(pollID uint, userID uint) (*Poll, bool)
    DeletePollVote(pollId uint, userId uint) error
    GetActivePoll() (*Poll, error)
}
