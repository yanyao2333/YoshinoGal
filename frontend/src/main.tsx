import * as React from "react";
import {HashRouter, Route, Routes} from "react-router-dom";
import LibraryPage from "./pages/libraryPage";

import './globals.css'
import {createRoot} from "react-dom/client";

const domNode = document.getElementById('root');
// @ts-ignore
const root = createRoot(domNode);

root.render(
    // <HashRouter>
    //     <Routes>
    //         {/*<Route path="/game/:id" element={<GamePage/>}/>*/}
    //         <Route path="/library" element={<LibraryPage/>}/>
    //     </Routes>
    // </HashRouter>,
    <HashRouter basename={"/"}>
        {/* The rest of your app goes here */}
        <Routes>
            <Route path="/library" element={<LibraryPage/>}/>
            <Route path="/" element={<LibraryPage/>}/>
        </Routes>
    </HashRouter>
);