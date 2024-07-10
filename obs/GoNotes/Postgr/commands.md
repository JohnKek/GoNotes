distinct - уникальные

count() - подсчет

between 20 and 30 попадение в диапазон между

in, not in (list)

order by: select distinct * from employees order by first_name

ASC DESC select distinct * from employees order by first_name,employee_id desc

Скалярные функции MAX MIN AVG

Pattern Matching LIKE
% хоть сколько символов
_ 1 символ

LIMIT ограничение выборки

Check NULL
SELECT ship_city,ship_region FROM orders where ship_city IS NULL

GROUP BY

SELECT ship_country, COUNT(*)
FROM orders
WHERE freight > 50
GROUP BY ship_country
ORDER BY COUNT(*) DESC 


HAVING
SELECT products.category_id, SUM(products.unit_price + products.units_in_stock)
FROM products
WHERE discontinued <> 1
GROUP BY products.category_id
HAVING  SUM(products.unit_price + products.units_in_stock) > 200


UNION
SELECT customers.country from customers union select employees.country from employees объединение результатов без дубликатов
UNION ALL c дубликатами
INTERSECT пересечение множеств
EXCEPT исключение 2 из 1 