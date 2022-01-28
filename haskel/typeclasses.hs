import qualified Data.Map as Map

data Point = Point Float Float deriving (Show)
data Shape = Circle Point Float | Rectangle Point Point deriving (Show)

surface :: Shape -> Float
surface (Circle _ r) = pi * r ^ 2
surface (Rectangle (Point x1 y1) (Point x2 y2)) = (abs $ x2 - x1) * (abs $ y2 - y1)

nudge :: Shape -> Float -> Float -> Shape
nudge (Circle (Point x y) r) a b = Circle (Point (x+a) (y+b)) r
nudge (Rectangle (Point x1 y1) (Point x2 y2)) a b = Rectangle (Point (x1+a) (y1+b)) (Point (x2+a) (y2+b))

data Person = Person { firstName :: String
                     , lastName :: String
                     , age :: Int
                     , height :: Float
                     , phoneNumber :: String
                     , flavor :: String
                     } deriving (Show)

data PersonShort = PersonShort { firstName' :: String
                               , lastName' :: String
                               , age' :: Int
                               } deriving (Eq, Show, Read)

data Car = Car { company :: String
               , model :: String
               , year :: Int
               } deriving (Show)

data Vector a = Vector a a a deriving (Show)

vplus :: Num a => Vector a -> Vector a -> Vector a
(Vector i j k) `vplus` (Vector l m n) = Vector (i + l) (j + m) (k + n) 

vectMult :: (Num t) => Vector t -> t -> Vector t  
(Vector i j k) `vectMult` m = Vector (i*m) (j*m) (k*m)  
  
scalarMult :: (Num t) => Vector t -> Vector t -> t  
(Vector i j k) `scalarMult` (Vector l m n) = i*l + j*m + k*n

data Day = Monday | Tuesday | Wednesday | Thursday | Friday | Saturday | Sunday   
           deriving (Eq, Ord, Show, Read, Bounded, Enum)


data LockerState = Taken | Free deriving (Show, Eq)
type Code = String
type LockerMap = Map.Map Int (LockerState, Code)

lockerLookup :: Int -> LockerMap -> Either String Code
lockerLookup lockerNumber map = 
    case Map.lookup lockerNumber map of
        Nothing -> Left $ "Locker number " ++ show lockerNumber ++ " does not exist!"
        Just (state, code) -> if state /= Taken
                                then Right code
                                else Left $ "Locker " ++ show lockerNumber ++ " is already taken!"


data Tree a = EmptyTree | Node a (Tree a) (Tree a) deriving (Show)

singleton :: a -> Tree a
singleton x = Node x EmptyTree EmptyTree

treeInsert :: (Ord a) => a -> Tree a -> Tree a
treeInsert x EmptyTree = singleton x
treeInsert x (Node a left right)
    | x == a = Node x left right
    | x < a  = Node a (treeInsert x left) right
    | x > a  = Node a left (treeInsert x right)

treeElem :: Ord a => a -> Tree a -> Bool
treeElem x EmptyTree = False 
treeElem x (Node a left right)
    | x == a = True 
    | x < a  = treeElem x left
    | x > a  = treeElem x right

data TrafficLight = Red | Yellow | Green

instance Eq TrafficLight where
    Red == Red = True 
    Yellow == Yellow = True 
    Green == Green = True 
    _ == _ = False 

instance Show TrafficLight where
    show Red = "Red light"
    show Yellow = "Yellow light"
    show Green = "Green light"

-- YesNo typeclass
class YesNo a where
    yesno :: a -> Bool

instance YesNo Int where
    yesno 0 = False
    yesno _ = True

instance YesNo [a] where
    yesno [] = False
    yesno _ = True

instance YesNo Bool where
    yesno = id

instance YesNo (Maybe a) where
    yesno (Just _) = True
    yesno Nothing = False

instance YesNo (Tree a) where
    yesno EmptyTree = False
    yesno _ = True

instance YesNo TrafficLight where
    yesno Red = False
    yesno _ = True

yesnoIf :: (YesNo y) => y -> a -> a -> a
yesnoIf yesnoVal yesResult noResult = if yesno yesnoVal then yesResult else noResult

-- functor
instance Functor Tree where
    fmap f EmptyTree = EmptyTree
    fmap f (Node x leftsub rightsub) = Node (f x) (fmap f leftsub) (fmap f rightsub)

main :: IO()
main = do
    print $ surface $ Rectangle (Point 0 0) (Point 100 100)
    print $ Circle (Point 10 20) 5
    print $ surface $ Circle (Point 10 20) 24
    print $ nudge (Circle (Point 34 34) 10) 5 10

    print $ map (Circle (Point 10 20)) [4,5,6,6]

    print $ Car { company="Ford", model="Mustang", year=1967}

    let mikeD = PersonShort {firstName' = "Michael", lastName' = "Diamond", age' = 43}
    let adRock = PersonShort {firstName' = "Adam", lastName' = "Horovitz", age' = 41}
    let mca = PersonShort {firstName' = "Adam", lastName' = "Yauch", age' = 44}

    print $ firstName' mikeD

    print $ mca == adRock
    print $ mikeD == mikeD
    print $ mikeD == PersonShort {firstName' = "Michael", lastName' = "Diamond", age' = 43}

    print mikeD

    let personShort = read "PersonShort {firstName' =\"Michael\", lastName' =\"Diamond\", age' = 43}" :: PersonShort
    print personShort

    print [Thursday .. Sunday]
    print ([minBound .. maxBound ] :: [Day])
 
    let lockers = Map.fromList [(100,(Taken,"ZD39I"))  
            ,(101,(Free,"JAH3I"))  
            ,(103,(Free,"IQSA9"))  
            ,(105,(Free,"QOTSA"))  
            ,(109,(Taken,"893JJ"))  
            ,(110,(Taken,"99292"))  
            ]  

    print $ lockerLookup 101 lockers
    print $ lockerLookup 100 lockers
    print $ lockerLookup 102 lockers
    print $ lockerLookup 110 lockers
    print $ lockerLookup 105 lockers

    -- recursive data structures 
    let nums = [8,6,4,1,7,3,5]
    let numsTree = foldr treeInsert EmptyTree nums
    print numsTree
    print $ 8 `treeElem` numsTree
    print $ 100 `treeElem` numsTree

    -- typeclasses 102
    print $ Red == Yellow
    print $ Red == Red
    print $ Red `elem` [Red, Yellow, Green]
    print $ [Red, Yellow, Green]

    -- yesno
    print "YesNo"
    print $ yesno $ length []
    print $ yesno "haha"
    print $ yesno $ Just 0
    print $ yesno [0,0,0]

    print $ yesnoIf [] "YEAH!" "NO!"
    print $ yesnoIf [2,3,4] "Yeah!" "No!"
    print $ yesnoIf (Just 500) "YEAH!" "No!"
    print $ yesnoIf Nothing "YEAH!" "NO!"

    -- functor
    print $ fmap (*2) EmptyTree
    print $ fmap (*4) (foldr treeInsert EmptyTree [5,7,3,2,1,7])
