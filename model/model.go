package model

// 导入包
import (
	"gorm.io/gorm"
)

// User 定义User结构体，对应user表
type User struct {
	ID              int64          `gorm:"primaryKey;autoIncrement;not null" json:"id"`  // 主键、自增、不为空
	Name            string         `gorm:"type:varchar(64);unique;not null" json:"name"` // 名称为唯一索引、不为空
	Password        string         `gorm:"type:varchar(255)" json:"-"`                   // 忽略密码字段
	FollowCount     int32          `gorm:"default:0;not null" json:"follow_count"`       // 默认0、不为空
	FollowerCount   int32          `gorm:"default:0;not null" json:"follower_count"`     // 默认0、不为空
	Avatar          string         `gorm:"type:varchar(255)" json:"avatar"`              // 用户头像地址
	BackgroundImage string         `gorm:"type:varchar(255)" json:"background_image"`    // 用户背景图地址
	Signature       string         `gorm:"type:varchar(255)" json:"signature"`           // 用户签名
	TotalFavorited  int32          `gorm:"default:0;not null" json:"total_favorited"`    // 忽略总收藏数字段
	WorkCount       int32          `gorm:"default:0;not null" json:"work_count"`         // 作品数量
	FavoriteCount   int32          `gorm:"default:0;not null" json:"favorite_count"`     // 收藏数
	CreateAt        int64          `gorm:"autoCreateTime" json:"-"`                      // 忽略创建时间字段
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`                               // 忽略删除标志字段
	IsFollow        bool           `gorm:"-" json:"is_follow"`                           //true-已关注，false-未关注

	Videos          []Video     `gorm:"foreignKey:AuthorID" json:"-"` // 忽略视频关联字段
	Comments        []Comment   `gorm:"foreignKey:UserID" json:"-"`   // 忽略评论关联字段
	Followers       []UserUser  `gorm:"foreignKey:Followed" json:"-"` // 忽略粉丝关联字段
	Followeds       []UserUser  `gorm:"foreignKey:Follower" json:"-"` // 忽略关注关联字段
	FavoritedVideos []UserVideo `gorm:"foreignKey:UserID" json:"-"`   // 忽略收藏视频关联字段
}

func (User) TableName() string {
	return "user"
}

// Video 定义Video结构体，对应video表
type Video struct {
	ID            int64          `gorm:"primaryKey;autoIncrement;not null" json:"id"` // 主键、自增、不为空
	AuthorID      int64          `gorm:"not null" json:"-"`                           // 不为空
	PlayURL       string         `gorm:"type:varchar(255);not null" json:"play_url"`  // 不为空
	CoverURL      string         `gorm:"type:varchar(255);not null" json:"cover_url"` // 不为空
	FavoriteCount int32          `gorm:"default:0;not null" json:"favorite_count"`    // 默认0、不为空
	CommentCount  int32          `gorm:"default:0;not null" json:"comment_count"`     // 默认0、不为空
	Title         string         `gorm:"type:varchar(255);not null" json:"title"`     // 不为空
	CreateAt      int64          `gorm:"autoCreateTime;index" json:"-"`               // 忽略创建时间字段
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`                              // 忽略删除标志字段                         // 0正常 1删除
	IsFavorite    bool           `gorm:"-" json:"is_favorite"`

	Author       *User      `gorm:"foreignKey:AuthorID" json:"author"` // 关联用户表，外键为AuthorID
	Comments     []Comment  `gorm:"foreignKey:VideoID" json:"-"`       // 关联评论表，外键为VideoID
	FavoriteUser *UserVideo `gorm:"foreignKey:VideoID" json:"-"`
}

func (Video) TableName() string {
	return "video"
}

// Comment 定义Comment结构体，对应comment表
type Comment struct {
	ID        int64          `gorm:"primaryKey;autoIncrement;not null" json:"id"` // 主键自增id、不为空
	UserID    int64          `gorm:"not null" json:"-"`                           // 外键，不为空
	VideoID   int64          `gorm:"index;not null" json:"-"`                     // 外键，索引，不为空
	Content   string         `gorm:"type:varchar(255);not null" json:"content"`   // 不为空
	CreateAt  int64          `gorm:"autoCreateTime" json:"-"`                     // 忽略创建时间字段
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                              // 0正常 1删除

	User *User `gorm:"foreignKey:UserID" json:"user"` // 关联用户表，外键为UserID
}

func (Comment) TableName() string {
	return "comment"
}

// UserUser 定义UserUser结构体，对应user_user表
type UserUser struct {
	ID       int64 `gorm:"primaryKey;autoIncrement;not null"` // 主键自增id、不为空
	Follower int64 `gorm:"uniqueIndex:er_ed_ui;not null"`     // 外键、联合唯一索引er_ed_ui、不为空
	Followed int64 `gorm:"uniqueIndex:er_ed_ui;not null"`     // 外键、联合唯一索引er_ed_ui、不为空
	Flag     bool  `gorm:"default:1;index;not null"`          // 字段为1表示关注 Follower 关注 Followed

	FollowerUser *User `gorm:"foreignKey:Follower" json:"-"` // 关联用户表，外键为Follower
	FollowedUser *User `gorm:"foreignKey:Followed" json:"-"` // 关联用户表，外键为Followed
}

func (UserUser) TableName() string {
	return "user_user"
}

// UserVideo 定义UserVideo结构体，对应user_video表
type UserVideo struct {
	ID      int64 `gorm:"primaryKey;autoIncrement;not null"`  // 主键自增id
	UserID  int64 `gorm:"uniqueIndex:user_video_ui;not null"` // 联合唯一索引user_video_ui、不为空
	VideoID int64 `gorm:"uniqueIndex:user_video_ui;not null"` // 联合唯一索引user_video_ui、不为空
	Flag    bool  `gorm:"default:1;index;not null" json:"-"`  // 字段为1表示 UserID 喜欢 VideoID

	User  *User  `gorm:"foreignKey:UserID" json:"-"`  // 关联用户表，外键为UserID
	Video *Video `gorm:"foreignKey:VideoID" json:"-"` // 关联视频表，外键为VideoID
}

func (UserVideo) TableName() string {
	return "user_video"
}
