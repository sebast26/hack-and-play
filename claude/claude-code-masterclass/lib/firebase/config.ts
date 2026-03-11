import { initializeApp } from "firebase/app";
import { getAuth } from "firebase/auth";
import { getFirestore } from "firebase/firestore";

const firebaseConfig = {
  apiKey: "AIzaSyDEwzb_tvUGqPTGdAu6_s3eI8ihnPCzYro",
  authDomain: "pocket-heist-website-sgorecki.firebaseapp.com",
  projectId: "pocket-heist-website-sgorecki",
  storageBucket: "pocket-heist-website-sgorecki.firebasestorage.app",
  messagingSenderId: "25934752507",
  appId: "1:25934752507:web:a42de741c9b5aaaa15de5e",
};

const app = initializeApp(firebaseConfig);

export const auth = getAuth(app);
export const db = getFirestore(app);

export default app;
