import { FieldValue } from "firebase/firestore";

export interface Heist {
  id: string;
  title: string;
  description: string;
  createdBy: string;
  createdByCodename: string;
  assignedTo: string;
  assignedToCodename: string;
  createdAt: Date;
  deadline: Date;
  finalStatus: "success" | "failure" | null;
}

export interface CreateHeistInput {
  title: string;
  description: string;
  createdBy: string;
  createdByCodename: string;
  assignedTo: string;
  assignedToCodename: string;
  createdAt: FieldValue;
  deadline: Date;
  finalStatus: null;
}

export interface UpdateHeistInput {
  title?: string;
  description?: string;
  assignedTo?: string;
  assignedToCodename?: string;
  deadline?: Date;
  finalStatus?: "success" | "failure" | null;
}
