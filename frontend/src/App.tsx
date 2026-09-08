import { Route, Routes } from "react-router-dom";
import { Layout } from "./components/Layout";
import { NotesDashboard } from "./pages/NotesDashboard";
import { CreateNotePage } from "./pages/CreateNotePage";
import { EditNotePage } from "./pages/EditNotePage";
import { NoteViewPage } from "./pages/NoteViewPage";

export default function App() {
  return (
    <Routes>
      <Route element={<Layout />}>
        <Route path="/" element={<NotesDashboard />} />
        <Route path="/notes/new" element={<CreateNotePage />} />
        <Route path="/notes/:id" element={<NoteViewPage />} />
        <Route path="/notes/:id/edit" element={<EditNotePage />} />
      </Route>
    </Routes>
  );
}
