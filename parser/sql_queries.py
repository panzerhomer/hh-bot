_CREATE_TABLE_QUERY = """
    CREATE TABLE IF NOT EXISTS vacancies (
        id SERIAL PRIMARY KEY,
        city VARCHAR(50),
        company VARCHAR(200),
        industry VARCHAR(200),
        title VARCHAR(200),
        keywords TEXT,
        skills  TEXT,
        experience VARCHAR(50),
        salary VARCHAR(50),
        url VARCHAR(200)
    )
"""


_DROP_TABLE_QUERY = "DROP TABLE IF EXISTS vacancies"


_DELETE_DUPLICATES_QUERY = """
            DELETE FROM vacancies
            WHERE id NOT IN (
                SELECT MIN(id)
                FROM vacancies
                GROUP BY url
            )
        """
