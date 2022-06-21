SELECT author_name, count(*)
FROM commits -- replace repo
WHERE parents < 2 -- ignore merge commits
GROUP BY author_name ORDER BY count(*) DESC
LIMIT 20

