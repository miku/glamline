-- top 50 files changed most frequently in the past year
SELECT file_path, COUNT(*)
FROM commits('/home/tir/code/rclone/rclone'), stats('/home/tir/code/rclone/rclone', commits.hash)
WHERE
commits.author_when > DATE('now', '-12 month')
AND commits.parents < 2 -- ignore merge commits
GROUP  BY file_path
ORDER  BY COUNT(*) DESC
LIMIT  20
