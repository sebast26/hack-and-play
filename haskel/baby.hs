doubleMe x = x * 2
doubleUs x y = doubleMe x + doubleMe y

doubleSmallNumber x = if x <= 100
                        then x * 2
                        else x

conanO'Brien = "It's a-me, Conan O'Brien!"

"hello" ++ " " ++ "world" 
5:[1,2,3,4,5]  
"Steve Buscemi" !! 6
head [5,4,3,2,1]
tail [5,4,3,2,1]
last [5,4,3,2,1]
init [5,4,3,2,1]
length [5,4,3,2,1]
null [1,2,3]
null [] 
reverse [5,4,3,2,1]
take 3 [5,4,3,2,1]
take 5 [1,2]  
drop 3 [8,4,2,1,5,6]
maximum [8,4,2,1,5,6]
minimum [8,4,2,1,5,6]
sum [8,4,2,1,5,6]
product [6,2,1,2]
4 `elem` [3,4,5,6]
10 `elem` [3,4,5,6]


[1..20]
['K'..'Z']
[2,4..20]
[3,6,20]

[20,19..1]

take 24 [13,26..]

take 10 (cycle [1,2,3])
take 12 (cycle "LOL ")
take 10 (repeat 5)
replicate 3 10


[x*2 | x <- [1..10]]
[x*2 | x <- [1..10], x*2 >= 12]
[ x | x <- [50..100], x `mod` 7 == 3]
boomBangs xs = [ if x < 10 then "BOOM!" else "BANG!" | x <- xs, odd x]

[ x*y | x <- [2,5,10], y <- [8,10,11], x*y > 50]



let xxs = [[1,3,5,2,3,1,2,4,5],[1,2,3,4,5,6,7,8,9],[1,2,4,2,1,6,3,1,3,2,3,6]]
[ [ x | x <- xs, even x ] | xs <- xxs]