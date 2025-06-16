package db

import (
	"context"
	"database/sql"
	"log"
	"math/rand/v2"
	"regexp"  // Added for sanitization
	"strings" // Added for sanitization

	"github.com/LincolnG4/Haku/internal/store"
)

// Cleaned up the ghibliCharacters slice slightly (removed trailing space from Chihiro).
var ghibliCharacters = []string{
	// --- Spirited Away (Mainly) ---
	"Chihiro Ogino", // Removed trailing space
	"Haku",
	"Yubaba",
	"Zeniba",
	"Kamaji",
	"Lin",
	"No-Face",
	"Boh",
	"Radish Spirit",
	"Aogaeru",
	"Kashira",
	"River Spirit",
	"Soot Sprites",
	"Akio Ogino",
	"Yūko Ogino",
	"Yu-bird",

	// --- Princess Mononoke ---
	"San (Princess Mononoke)",
	"Ashitaka",
	"Lady Eboshi",
	"Jigo",
	"Moro (Wolf God)",
	"Yakul (Red Elk)",
	"Kodama (Tree Spirits)",

	// --- Howl's Moving Castle ---
	"Sophie Hatter",
	"Howl Jenkins Pendragon",
	"Calcifer",
	"Markl",
	"Witch of the Waste",
	"Turnip Head (Prince Justin)",

	// --- My Neighbor Totoro ---
	"Satsuki Kusakabe",
	"Mei Kusakabe",
	"Totoro",
	"Catbus (Nekobasu)",
	"Kanta Ōgaki",
	"Granny",

	// --- Kiki's Delivery Service ---
	"Kiki",
	"Jiji",
	"Tombo Kopoli",
	"Ursula",
	"Osono",

	// --- Ponyo ---
	"Ponyo (Brunhilde)",
	"Sosuke",
	"Lisa",
	"Fujimoto",
	"Granmamare",

	// --- Other Classics ---
	"Sheeta (Castle in the Sky)",
	"Pazu (Castle in the Sky)",
	"Colonel Muska (Castle in the Sky)",
	"Nausicaä (Nausicaä of the Valley of the Wind)",
	"Porco Rosso (Porco Rosso)",
}

var emailInvalidChars = regexp.MustCompile(`[^a-z0-9.]+`)

func formatEmailLocalPart(name string) string {
	// 1. Convert to lowercase
	s := strings.ToLower(name)

	// 2. Replace spaces with dots
	s = strings.ReplaceAll(s, " ", ".")

	// 3. Remove all characters that are not lowercase letters, numbers, or dots
	s = emailInvalidChars.ReplaceAllString(s, "")

	return s
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(50)
	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Printf("error seeding user %s: %v", user.Username, err)
			return
		}
	}

	pipelines := generatePipelines(10, users)
	for _, pipeline := range pipelines {
		if err := store.Pipelines.Create(ctx, pipeline); err != nil {
			log.Printf("error seeding pipeline %s: %v", pipeline.Name, err)
			return
		}
	}

	log.Println("seeding completed")
}

func generateUsers(n int) []*store.User {
	if n > len(ghibliCharacters) {
		n = len(ghibliCharacters)
	}
	users := make([]*store.User, n)

	for i := 0; i < n; i++ {
		name := ghibliCharacters[i]

		emailLocalPart := formatEmailLocalPart(name)

		users[i] = &store.User{
			Username: name,
			Email:    emailLocalPart + "@ghibli.com",
		}
	}
	return users
}

func generatePipelines(n int, users []*store.User) []*store.Pipelines {
	pipelines := make([]*store.Pipelines, n)
	numUsers := len(users)

	if numUsers == 0 {
		return pipelines
	}

	for i := 0; i < n; i++ {
		user := users[rand.IntN(numUsers)]
		pipelines[i] = &store.Pipelines{
			UserID: user.ID,
			Name:   user.Username + "'s pipeline",
		}
	}
	return pipelines
}
