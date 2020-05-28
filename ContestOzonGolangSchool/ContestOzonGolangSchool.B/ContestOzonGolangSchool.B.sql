SELECT g.id, g.name
FROM tags_goods tg
JOIN goods g ON g.id = tg.goods_id
GROUP BY g.id, g.name
HAVING COUNT(1) = (SELECT COUNT(1) FROM tags)