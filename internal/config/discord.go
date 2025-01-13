package config

import "os"
import _ "github.com/joho/godotenv/autoload"

type Rating struct {
	ShortName     string
	CertValue     int
	DisocrdRoleId string
}

type Ratings struct {
	Ratings []Rating
}

var DiscordToken = os.Getenv("DISCORD_TOKEN")

func GetRatingsRoles() Ratings {

	return Ratings{
		Ratings: []Rating{
			{
				ShortName:     "ADM",
				CertValue:     12,
				DisocrdRoleId: "1037910657779126292",
			},
			{
				ShortName:     "SUP",
				CertValue:     11,
				DisocrdRoleId: "1037908270758768651",
			},
		},
	}
}

func GetPilotRatingRoles() Ratings {
	return Ratings{
		Ratings: []Rating{
			{
				ShortName:     "P0",
				CertValue:     0,
				DisocrdRoleId: "1037910884590301226",
			},
			{
				ShortName:     "P1",
				CertValue:     1,
				DisocrdRoleId: "1037910919528841257",
			},
			{
				ShortName:     "P2",
				CertValue:     3,
				DisocrdRoleId: "1037911208650608670",
			},
			{
				ShortName:     "P3",
				CertValue:     7,
				DisocrdRoleId: "1037911259150032987",
			},
			{
				ShortName:     "P4",
				CertValue:     15,
				DisocrdRoleId: "1037914612407992330",
			},
			{
				ShortName:     "P5",
				CertValue:     31,
				DisocrdRoleId: "1220953753105203291",
			},
			{
				ShortName:     "P6",
				CertValue:     63,
				DisocrdRoleId: "1037908270758768658",
			},
		},
	}
}
