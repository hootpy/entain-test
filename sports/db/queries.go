package db

const (
	eventsList = "list"
)

func getEventQueries() map[string]string {
	return map[string]string{
		eventsList: `
			WITH e AS (
				SELECT
					id,
					name,
					advertised_start_time,
					visible,
					CASE
						WHEN datetime(advertised_start_time) <= datetime(?) THEN 'CLOSED'
						ELSE 'OPEN'
					END AS status
				FROM events
			)
			SELECT *
			FROM e
		`,
	}
}
