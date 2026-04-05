package db

const (
	racesList = "list"
)

func getRaceQueries() map[string]string {
	return map[string]string{
		racesList: `
			WITH r AS (
				SELECT
					id,
					meeting_id,
					name,
					number,
					visible,
					advertised_start_time,
					CASE
						WHEN datetime(advertised_start_time) <= datetime(?) THEN 'CLOSED'
						ELSE 'OPEN'
					END AS status
				FROM races
			)
			SELECT *
			FROM r
		`,
	}
}
