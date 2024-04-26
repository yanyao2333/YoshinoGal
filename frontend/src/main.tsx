import * as React from "react";
import * as ReactDOM from "react-dom";
import {HashRouter, Route, Routes} from "react-router-dom";
import LibraryPage from "./pages/libraryPage";

import './globals.css'

ReactDOM.render(
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
    </HashRouter>,
    document.getElementById('root')
);