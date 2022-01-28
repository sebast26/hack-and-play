module Main where
import qualified Data.List as L
import qualified Data.Function as F
import qualified Data.Char as C
import qualified Data.Map as Map
import qualified Data.Set as Set

numUniques :: (Ord a) => [a] -> Int
numUniques = length . L.nub

numUniques' :: (Ord a) => [a] -> Int
numUniques' = \xs -> length (L.nub xs)

main :: IO()
main = do
    let a = L.intersperse '.' "MONKEY" -- przeplatać
    print a -- "M.O.N.K.E.Y"
    let b = L.intercalate [0,0,0] [[1,2,3],[4,5,6],[7,8,9]] 
    print b -- [1,2,3,0,0,0,4,5,6,0,0,0,7,8,9] 
    let c = L.transpose [[1,2,3],[4,5,6],[7,8,9]] -- 2D matrix, the columns become the rows and vice versa
    print c -- [[1,4,7],[2,5,8],[3,6,9]] 
    let d = L.concat ["foo","bar","car"]
    print d -- "foobarcar"
    let e = L.concat [[3,4,5],[2,3,4],[2,1,1]]
    print e -- [3,4,5,2,3,4,2,1,1]
    let f = L.concatMap (replicate 4) [1..3] 
    print f -- [1,1,1,1,2,2,2,2,3,3,3,3]
    let g = and $ map (>4) [5,6,7,8]
    print g -- True
    let h = and $ map (==4) [4,4,4,3,4]
    print h -- False
    let i = any (==4) [2,3,5,6,1,4] 
    print i -- True
    let j = all (`elem` ['A'..'Z']) "HEYGUYSwhatsup" 
    print j -- False
    let k = take 3 $ iterate (++ "haha") "haha"
    print k -- ["haha","hahahaha","hahahahahaha"]
    let l = splitAt 3 "heyman"
    print l
    let m = takeWhile (/=' ') "This is a sentence" 
    print m -- "This"
    let n = sum $ takeWhile (<10000) $ map (^3) [1..]
    print n -- 53361
    let o = dropWhile (/=' ') "This is a sentence"
    print o -- " is a sentence"
    let stock = [(994.4,2008,9,1),(995.2,2008,9,2),(999.2,2008,9,3),(1001.4,2008,9,4),(998.3,2008,9,5)] 
    let p = head $ dropWhile (\(val, y, m , d) -> val < 1000) stock
    print p -- (1001.4,2008,9,4)
    let r = span (/=' ') "This is a sentence"
    print r -- ("This"," is a sentence")
    let s = break (==' ') "This is a sentence"
    print s -- ("This"," is a sentence")
    let t = L.sort [8,5,3,2,1,6,4,2]
    print t 
    let u = L.group [1,1,1,1,2,2,2,2,3,3,2,2,2,5,6,7]
    print u -- [[1,1,1,1],[2,2,2,2],[3,3],[2,2,2],[5],[6],[7]]
    -- If we sort a list before grouping it, we can find out how many times each element appears in the list.
    let w = map (\l@(x:xs) -> (x,length l)) . L.group . L.sort $ [1,1,1,1,2,2,2,2,3,3,2,2,2,5,6,7]
    print w -- [(1,4),(2,7),(3,2),(5,1),(6,1),(7,1)]
    -- searches for a sublist
    let x = "cat" `L.isInfixOf` "im a cat burglar"
    print x -- True
    print $ "hey" `L.isPrefixOf` "hey there!"
    print $ "there!" `L.isSuffixOf` "oh hey there!"
    let y = L.partition (>3) [1,3,5,6,3,2,1,0,3,7]
    print y -- ([5,6,7],[1,3,3,2,1,0,3])
    
    let z = L.find (>9) [1,2,3,4,5,6]
    print z -- Nothing

    let aa = L.find (\(val,y,m,d) -> val > 1000) stock
    print aa -- Just (1001.4,2008,9,4)
    let ab = 4 `L.elemIndex` [1,2,3,4,5,6]
    print ab -- Just 3
    let ac = ' ' `L.elemIndices` "Where are the spaces?"
    print ac -- [5,9,13]
    let ad = lines "first line\nsecond line\nthird line"
    print ad -- ["first line","second line","third line"]
    let ae = unlines ["first line", "second line", "third line"]
    print ae -- "first line\nsecond line\nthird line\n"
    let af = words "hey these are the words in this sentence"
    print af -- ["hey","these","are","the","words","in","this","sentence"]
    -- removes duplicates
    let ag = L.nub "Lots of words and stuff"
    print ag
    let ah = L.delete 'h' . L.delete 'h' . L.delete 'h' $ "hey there ghang!"
    print ah -- "ey tere gang!"
    -- list difference
    -- let ai = [1..10] `\\` [2,5,9]
    -- print ai -- [1,3,4,6,7,8,10]
    let aj = [1..7] `L.union` [5..10]
    print aj -- [1,2,3,4,5,6,7,8,9,10]
    let ak = [1..7] `L.intersect` [5..10]
    print ak
    let al = L.insert 4 [3,5,1,2,8,2]
    print al -- [3,4,5,1,2,8,2]

    let values = [-4.3, -2.4, -1.2, 0.4, 2.3, 5.9, 10.5, 29.1, 5.3, -2.4, -14.5, 2.9, 2.3]
    let am = L.groupBy (\x y -> (x > 0) == (y > 0)) values
    print am -- [[-4.3,-2.4,-1.2],[0.4,2.3,5.9,10.5,29.1,5.3],[-2.4,-14.5],[2.9,2.3]]
    let an = L.groupBy ((==) `F.on` (> 0)) values
    print an -- [[-4.3,-2.4,-1.2],[0.4,2.3,5.9,10.5,29.1,5.3],[-2.4,-14.5],[2.9,2.3]]

    let xs = [[5,4,5,4,4],[1,2,3],[3,5,4,3],[],[2],[2,2]]
    -- compare `on` length is equovalent of \x y -> length x `compare` length y
    let ao = L.sortBy (compare `F.on` length) xs
    print ao -- [[],[2],[2,2],[1,2,3],[3,5,4,3],[5,4,5,4,4]]

    print $ show $ C.isAlpha 'a'
    print $ show $ C.isDigit '4'
    print $ show $ C.isHexDigit 'A'

    print $ L.all C.isAlphaNum "bobby283"

    print $ show $ C.generalCategory ' ' -- Space
    print $ show $ C.generalCategory 'A' -- UppercaseLetter
    print $ map C.generalCategory " \t\nA9?|" -- [Space,Control,Control,UppercaseLetter,DecimalNumber,OtherPunctuation,MathSymbol]
    print $ C.generalCategory ' ' == C.Space -- True

    print $ map C.digitToInt "FF85AB" -- [15,15,8,5,10,11]
    
    let map = Map.fromList [("betty","555-2938"),("bonnie","452-2928"),("lucille","205-2928")]
    print map

    print $ Map.insert "sebastian" "333-1234" map

    print $ Map.lookup "betty" map

    print $ Map.map (*100) $ Map.fromList [(1,1), (2,4), (3,9)]

    print $ Map.filter C.isUpper $ Map.fromList [(1,'a'), (2, 'A'), (3, 'b'), (4, 'B')]

    let phoneBook =  [("betty","555-2938")  
            ,("betty","342-2492")  
            ,("bonnie","452-2928")  
            ,("patsy","493-2928")  
            ,("patsy","943-2929")  
            ,("patsy","827-9162")  
            ,("lucille","205-2928")  
            ,("wendy","939-8282")  
            ,("penny","853-2492")  
            ,("penny","555-2111")  
            ]

    print $ Map.fromListWith (\number1 number2 -> number1 ++ ", " ++ number2) phoneBook

    print $ Map.fromListWith max [(2,3),(2,5),(2,100),(3,29),(3,22),(3,11),(4,22),(4,15)]

    print $ Map.insertWith (+) 3 100 $ Map.fromList [(3,4),(5,103),(6,339)]

    let text1 = "I just had an anime dream. Anime... Reality... Are they so different?"
    let text2 = "The old man left his garbage can out and now his trash is all over my lawn!"
    let set1 = Set.fromList text1
    let set2 = Set.fromList text2
    print set1
    print set2

    print $ Set.intersection set1 set2
    print $ Set.difference set1 set2
    print $ Set.union set1 set2

    print $ Set.fromList [2,3,4] `Set.isSubsetOf` Set.fromList [1,2,3,4,5]