import { createUserWithEmailAndPassword, updateProfile } from "firebase/auth";
import { doc, setDoc } from "firebase/firestore";
import { auth, db } from "@/lib/firebase/config";
import { generateCodename } from "@/lib/codename";

export async function signUpUser(
  email: string,
  password: string,
): Promise<void> {
  const { user } = await createUserWithEmailAndPassword(auth, email, password);
  const codename = generateCodename();
  await updateProfile(user, { displayName: codename });
  await setDoc(doc(db, "users", user.uid), { id: user.uid, codename });
}
