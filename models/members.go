package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OpMembers struct {
	ID                primitive.ObjectID `bson:"_id"`
	Account           string             `bson:"account"`
	Password          string             `bson:"password"`
	Viplevel          string             `bson:"viplevel"`
	VipUpgradeAt      int64              `bson:"vip_upgrade_at"`
	PaymentGroup      string             `bson:"payment_group"`
	RegisterIP        string             `bson:"register_ip"`
	RegisterWebsite   string             `bson:"register_website"`
	AgentCode         string             `bson:"agent_code"`
	Recommender       string             `bson:"recommender"`
	RecommendCode     string             `bson:"recommend_code"`
	Status            string             `bson:"status"`
	LastLoginIP       string             `bson:"last_login_ip"`
	RebateRequirement int64              `bson:"rebate_requirement"`
	Birthdate         string             `bson:"birthdate"`
	CompletedBaseinfo string             `bson:"completed_baseinfo"`
	ContactNo         string             `bson:"contact_no"`
	Email             string             `bson:"email"`
	Gender            string             `bson:"gender"`
	Username          string             `bson:"username"`
	WithdrawalKey     string             `bson:"withdrawal_key"`
	VipChangedAt      int64              `bson:"vip_changed_at"`
	ChargeAmount      int64              `bson:"charge_amount"`
	ChargeCount       int64              `bson:"charge_count"`
	ChargeAt          int64              `bson:"charge_at"`
	Avatar            string             `bson:"avatar"`
	MemberCode        string             `bson:"member_code"`
	DepositTotal      int64              `bson:"deposit_total"`
	CreatedAt         time.Time          `bson:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at"`
	LineID            string             `bson:"line_id"`
	DisplayName       string             `bson:"display_name"`
	StatusMessage     string             `bson:"status_message"`
	PictureURL        string             `bson:"picture_url"`
	TransferPassword  string             `bson:"transfer_password"`
	WalletType        string             `bson:"wallet_type"`
}
