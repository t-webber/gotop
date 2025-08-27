package main

import (
	"database/sql"
	"log"
)

const QUERY = `
SELECT SUM(end - start), cwd
FROM processes
WHERE cwd IS NOT NULL
  AND cmdline LIKE '%nvim %'
  AND cmdline NOT LIKE '% --embed %'
GROUP BY cwd
ORDER BY SUM(end - start) DESC;
`

func getNvimUsage(db *sql.DB) []weightedData {
	rows, err := db.Query(QUERY)
	if err != nil {
		log.Fatalf("Failed to select nvim processes from database: %s", err)
	}

	projects := []weightedData{}

	for rows.Next() {
		var project weightedData
		if err := rows.Scan(&project.weight, &project.data); err != nil {
			log.Fatal(err)
		}
		projects = append(projects, project)

	}

	rows.Close()

	return projects
}
