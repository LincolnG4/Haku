package db

import (
	"context"
	"database/sql"
	"log"
	"math/rand/v2"
	"regexp"  // Added for sanitization
	"strings" // Added for sanitization

	storepkg "github.com/LincolnG4/Haku/internal/store"
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

// Mapping from movie to characters
var movieCharacters = map[string][]string{
	"Spirited Away": {
		"Chihiro Ogino", "Haku", "Yubaba", "Zeniba", "Kamaji", "Lin", "No-Face", "Boh", "Radish Spirit", "Aogaeru", "Kashira", "River Spirit", "Soot Sprites", "Akio Ogino", "Yūko Ogino", "Yu-bird",
	},
	"Princess Mononoke": {
		"San (Princess Mononoke)", "Ashitaka", "Lady Eboshi", "Jigo", "Moro (Wolf God)", "Yakul (Red Elk)", "Kodama (Tree Spirits)",
	},
	"Howl's Moving Castle": {
		"Sophie Hatter", "Howl Jenkins Pendragon", "Calcifer", "Markl", "Witch of the Waste", "Turnip Head (Prince Justin)",
	},
	"My Neighbor Totoro": {
		"Satsuki Kusakabe", "Mei Kusakabe", "Totoro", "Catbus (Nekobasu)", "Kanta Ōgaki", "Granny",
	},
	"Kiki's Delivery Service": {
		"Kiki", "Jiji", "Tombo Kopoli", "Ursula", "Osono",
	},
	"Ponyo": {
		"Ponyo (Brunhilde)", "Sosuke", "Lisa", "Fujimoto", "Granmamare",
	},
	"Castle in the Sky": {
		"Sheeta (Castle in the Sky)", "Pazu (Castle in the Sky)", "Colonel Muska (Castle in the Sky)",
	},
	"Nausicaä of the Valley of the Wind": {
		"Nausicaä (Nausicaä of the Valley of the Wind)",
	},
	"Porco Rosso": {
		"Porco Rosso (Porco Rosso)",
	},
}

// Reverse mapping: character -> movie
var characterMovie = func() map[string]string {
	m := make(map[string]string)
	for movie, chars := range movieCharacters {
		for _, c := range chars {
			m[c] = movie
		}
	}
	return m
}()

func Seed(store storepkg.Storage, db *sql.DB) {
	ctx := context.Background()

	// 1. Create organizations (movies)
	orgIDs := make(map[string]int64)
	for movie := range movieCharacters {
		org := &storepkg.Organization{
			Name:        movie,
			Description: "Organization for the movie " + movie,
		}
		if err := store.Organizations.Create(ctx, org); err != nil {
			log.Printf("error seeding organization %s: %v", org.Name, err)
			return
		}
		orgIDs[movie] = org.ID
	}

	// 2. Create users and assign to organizations
	users := generateUsersWithOrgs(ghibliCharacters, orgIDs)
	for _, user := range users {
		// Set a default password for all users
		_ = user.Password.Set("password123")
		if err := store.Users.Create(ctx, &user.User); err != nil {
			if strings.Contains(err.Error(), "duplicate key value") {
				log.Printf("user %s already exists, skipping", user.Username)
				continue
			}
			log.Printf("error seeding user %s: %v", user.Username, err)
			return
		}
		// Add user to organization as Admin
		member := &storepkg.OrganizationMember{
			UserID:         user.ID, // user.ID is set after Create
			OrganizationID: user.OrganizationID,
			RoleID:         storepkg.AdminRole,
		}
		if err := store.Organizations.AddMember(ctx, member); err != nil {
			log.Printf("error adding user %s to org: %v", user.Username, err)
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

// Helper struct to keep org info with user

type userWithOrg struct {
	storepkg.User
	OrganizationID int64
}

func generateUsersWithOrgs(characters []string, orgIDs map[string]int64) []userWithOrg {
	users := make([]userWithOrg, len(characters))
	for i, name := range characters {
		emailLocalPart := formatEmailLocalPart(name)
		movie := characterMovie[name]
		orgID := orgIDs[movie]
		users[i] = userWithOrg{
			User: storepkg.User{
				Username: name,
				Email:    emailLocalPart + "@ghibli.com",
			},
			OrganizationID: orgID,
		}
	}
	return users
}

func generatePipelines(n int, users []userWithOrg) []*storepkg.Pipelines {
	pipelines := make([]*storepkg.Pipelines, n)
	numUsers := len(users)

	if numUsers == 0 {
		return pipelines
	}

	for i := 0; i < n; i++ {
		user := users[rand.IntN(numUsers)]
		pipelines[i] = &storepkg.Pipelines{
			OrganizationID: user.OrganizationID,
			Name:           user.Username + "'s pipeline",
		}
	}
	return pipelines
}
