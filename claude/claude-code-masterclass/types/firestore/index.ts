import { DocumentData, QueryDocumentSnapshot } from "firebase/firestore";
import { Heist } from "./heist";

export * from "./heist";

export const COLLECTIONS = {
  HEISTS: "heists",
} as const;

export const heistConverter = {
  toFirestore: (data: Partial<Heist>): DocumentData => data,

  fromFirestore: (snapshot: QueryDocumentSnapshot): Heist =>
    ({
      id: snapshot.id,
      ...snapshot.data(),
      createdAt: snapshot.data().createdAt?.toDate(),
      deadline: snapshot.data().deadline?.toDate(),
    }) as Heist,
};
