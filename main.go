package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// enlever les accents dans la premier ligne du fichier txt
func main() {
	file, err := os.Open("resultats-par-niveau-burvot-t1-france-entiere.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	totalVotes := 0

	if scanner.Scan() {
		headers := strings.Split(scanner.Text(), ";")
		votersIndex := findColumnIndex(headers, "Votants")
		nomCandidatIndex := findColumnIndex(headers, "Nom")
		prenomCandidatIndex := findColumnIndex(headers, "Prenom")
		voixCandidatIndex := findColumnIndex(headers, "Voix")
		departementIndex := findColumnIndex(headers, "Libelle du departement")

		if votersIndex == -1 || prenomCandidatIndex == -1 || voixCandidatIndex == -1 || nomCandidatIndex == -1 || departementIndex == -1 {
			log.Fatal("Probleme sur une colonnne")
		}

		totalVotesParCandidat := make(map[string]int)
		totalVotesParDepartement := make(map[string]map[string]int)

		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Split(line, ";")

			if departementIndex < len(fields) && votersIndex < len(fields) {
				departement := fields[departementIndex]

				voters, err := strconv.Atoi(fields[votersIndex])
				if err != nil {
					log.Printf("Erreur lors de la conversion du nombre de votants : %v\n", err)
					continue
				}

				nomCandidat := fields[nomCandidatIndex]
				prenomCandidat := fields[prenomCandidatIndex]
				candidateKey := fmt.Sprintf("%s %s", nomCandidat, prenomCandidat)

				for i := voixCandidatIndex + 5; i < len(fields); i += 7 {
					nouveauNom := fields[i]
					nouveauPrenom := fields[i+1]
					voixCandidat, err := strconv.Atoi(fields[i+2])
					if err != nil {
						continue
					}
					nouvelleCleCandidat := fmt.Sprintf("%s %s", nouveauNom, nouveauPrenom)
					totalVotesParCandidat[nouvelleCleCandidat] += voixCandidat

					if totalVotesParDepartement[departement] == nil {
						totalVotesParDepartement[departement] = make(map[string]int)
					}
					totalVotesParDepartement[departement][nouvelleCleCandidat] += voixCandidat
				}

				voixCandidat, err := strconv.Atoi(fields[voixCandidatIndex])
				if err != nil {
					log.Printf("Erreur lors de la conversion du nombre de votants : %v\n", err)
					continue
				}
				totalVotesParCandidat[candidateKey] += voixCandidat

				if totalVotesParDepartement[departement] == nil {
					totalVotesParDepartement[departement] = make(map[string]int)
				}
				totalVotesParDepartement[departement][candidateKey] += voixCandidat

				totalVotes += voters
			}
		}

		// total des votes pour chaque candidat
		for candidate, totalVotes := range totalVotesParCandidat {
			fmt.Printf("Total des votes pour le candidat %s : %d\n", candidate, totalVotes)
		}

		// total des votants
		fmt.Printf("Total des votants : %d\n", totalVotes)

		// total des votes par département et par candidat
		for departement, votesParCandidat := range totalVotesParDepartement {
			fmt.Printf("Département de %s\n", departement+" : ")
			for candidat, votes := range votesParCandidat {
				fmt.Printf("\t%s : %d\n", candidat, votes)
			}
		}
		departements := make([]struct {
			Nom   string
			Votes int
		}, 0)

		// total des votes par département
		for departement, votesParCandidat := range totalVotesParDepartement {
			totalVotes := 0
			for _, votes := range votesParCandidat {
				totalVotes += votes
			}

			departements = append(departements, struct {
				Nom   string
				Votes int
			}{Nom: departement, Votes: totalVotes})
		}
		for i := 0; i < len(departements)-1; i++ {
			for j := i + 1; j < len(departements); j++ {
				if departements[i].Votes < departements[j].Votes {
					departements[i], departements[j] = departements[j], departements[i]
				}
			}
		}

		//Palmares des départements par nombre de votant
		fmt.Println("\nDépartements avec le nombre de votants du plus grand au plus petit:")
		for _, dep := range departements {
			fmt.Printf("%s : %d votants\n", dep.Nom, dep.Votes)
		}

	} else {
		log.Fatal("Le fichier est vide.")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func findColumnIndex(headers []string, columnName string) int {
	for i, header := range headers {
		if header == columnName {
			return i
		}
	}
	return -1
}
