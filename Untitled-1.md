curl 'https://leetcode.com/graphql/' -X POST -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/113.0' -H 'Accept: */*' -H 'Accept-Language: en-US,en;q=0.5' -H 'Accept-Encoding: gzip, deflate, br' -H 'content-type: application/json' -H 'x-csrftoken: KEVQh6oDZn1cF74aZuHpGe3bDSnajMjBlavDKX0hoBQo6AV89k1vrJMpUlwHYxqk' -H 'authorization: ' -H 'random-uuid: b9340b14-b459-0830-a29d-50ff9b8ec835' -H 'sentry-trace: bee13a75bca74b1e902334e3f23d53ca-abfcbbdd6948f2a9-0' -H 'baggage: sentry-environment=production,sentry-release=c40975a6,sentry-transaction=%2Fu%2F%5Busername%5D,sentry-public_key=2a051f9838e2450fbdd5a77eb62cc83c,sentry-trace_id=bee13a75bca74b1e902334e3f23d53ca,sentry-sample_rate=0.004' -H 'Origin: https://leetcode.com' -H 'Connection: keep-alive' -H 'Referer: https://leetcode.com/nairvarun/' -H 'Cookie: csrftoken=KEVQh6oDZn1cF74aZuHpGe3bDSnajMjBlavDKX0hoBQo6AV89k1vrJMpUlwHYxqk; LEETCODE_SESSION=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJfYXV0aF91c2VyX2lkIjoiNjYzODYzMyIsIl9hdXRoX3VzZXJfYmFja2VuZCI6ImRqYW5nby5jb250cmliLmF1dGguYmFja2VuZHMuTW9kZWxCYWNrZW5kIiwiX2F1dGhfdXNlcl9oYXNoIjoiYjY0Mjk4NjFlMzRmMDA2NDk0OGYwZTE4ZGFlN2YzOWFkMzgzNjFiNSIsImlkIjo2NjM4NjMzLCJlbWFpbCI6Im5haXJ2YXJ1bjEwNEBnbWFpbC5jb20iLCJ1c2VybmFtZSI6Im5haXJ2YXJ1biIsInVzZXJfc2x1ZyI6Im5haXJ2YXJ1biIsImF2YXRhciI6Imh0dHBzOi8vYXNzZXRzLmxlZXRjb2RlLmNvbS91c2Vycy9hdmF0YXJzL2F2YXRhcl8xNjc2NDM4OTUzLnBuZyIsInJlZnJlc2hlZF9hdCI6MTY4NDQyMjIxOSwiaXAiOiIxODIuNzIuMzkuMTAiLCJpZGVudGl0eSI6ImM5ZDIyZDVjNjUwMzMyNmZhYjk1OWI0ZTNkMzhiZGMwIiwic2Vzc2lvbl9pZCI6Mzg3OTg2OTMsIl9zZXNzaW9uX2V4cGlyeSI6MTIwOTYwMH0.0LDpLF-o1bqP4mITfy-YvENl_i1eGqiq5tgBsJOgZLQ; __stripe_mid=0bcf282e-1cff-4d76-aa74-405366f37ab7bca28c; _dd_s=rum=0&expire=1684425084650' -H 'Sec-Fetch-Dest: empty' -H 'Sec-Fetch-Mode: cors' -H 'Sec-Fetch-Site: same-origin' -H 'TE: trailers' --data-raw '{"query":"\n    query userProfileCalendar($username: String!, $year: Int) {\n  matchedUser(username: $username) {\n    userCalendar(year: $year) {\n      activeYears\n      streak\n      totalActiveDays\n      dccBadges {\n        timestamp\n        badge {\n          name\n          icon\n        }\n      }\n      submissionCalendar\n    }\n  }\n}\n    ","variables":{"username":"nairvarun"},"operationName":"userProfileCalendar"}'


{"1673740800": 5, "1673827200": 20, "1676332800": 2, "1676419200": 7, "1676505600": 12, "1676592000": 19, "1676678400": 16, "1676764800": 7, "1676851200": 1, "1676937600": 6, "1677024000": 1, "1677110400": 5, "1677196800": 1, "1677283200": 1, "1677369600": 1, "1677456000": 1, "1677542400": 5, "1677628800": 7, "1677715200": 11, "1677801600": 3, "1677888000": 2, "1677974400": 1, "1678060800": 1, "1678147200": 1, "1678233600": 1, "1679961600": 1, "1682985600": 32, "1683072000": 9, "1683676800": 10, "1683763200": 3, "1683849600": 20, "1683936000": 22, "1684022400": 2, "1684108800": 26, "1684281600": 11, "1684368000": 8, "1664668800": 8, "1669161600": 17, "1669248000": 1, "1669766400": 1}



1673740800: 5
1673827200: 20
1676332800: 2
1676419200: 7
1676505600: 12
1676592000: 19
1676678400: 16
1676764800: 7
1676851200: 1
1676937600: 6
1677024000: 1
1677110400: 5
1677196800: 1
1677283200: 1
1677369600: 1
1677456000: 1
1677542400: 5
1677628800: 7
1677715200: 11
1677801600: 3
1677888000: 2
1677974400: 1
1678060800: 1
1678147200: 1
1678233600: 1
1679961600: 1
1682985600: 32
1683072000: 9
1683676800: 10
1683763200: 3
1683849600: 20
1683936000: 22
1684022400: 2
1684108800: 26
1684281600: 11
1684368000: 8
1664668800: 8
1669161600: 17
1669248000: 1
1669766400: 1
1684432607
