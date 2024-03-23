import * as React from "react";
import * as ReactDOM from "react-dom";
import {HashRouter, Route, Routes} from "react-router-dom";
import LibraryPage from "./pages/libraryPage";

import './globals.css'
import GamePage from "./pages/gamePage";

ReactDOM.render(
    <HashRouter>
        <Routes>
            <Route path="/game/:id" element={<GamePage/>}/>
            <Route path="/library" element={<LibraryPage/>}/>
        </Routes>
    </HashRouter>,
    document.getElementById("root")
);