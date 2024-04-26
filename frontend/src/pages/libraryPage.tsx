import SideBar from "../components/sideBar";
import GamePosterWall from "../components/gamePosterWall";
import {GetPosterWall} from "../../wailsjs/go/app/Library";
import {useEffect, useState} from "react";

interface Game {
    name: string;
    imageSrc: string;
    id: number;
    href: string;
    imageAlt: string
}

async function getPosterWall() {
    let result = await GetPosterWall()
    let gamesShowList: Game[] = [];
    for (let game in result) {
        gamesShowList.push({
            name: result[game].game_name,
            imageSrc: result[game].poster_path,
            id: result[game].game_id,
            href: '/game/' + result[game].game_id,
            imageAlt: result[game].game_name,
        });
    }

    return gamesShowList
}


export default function Home() {
    const [games, setGames] = useState<Game[]>([]);

    useEffect(() => {
        getPosterWall().then(result => {
            setGames(result);
        }).catch(error => {
            console.error("Error fetching poster wall:", error);
        });
    }, []);

    return (
        <div className="flex bg-gray-100">
            <div className="flex flex-col gap-y-5 w-72 h-screen overflow-y-auto">
                <SideBar/>
            </div>
            <div className="flex-grow overflow-y-auto h-screen">
                <GamePosterWall game={games}/>
            </div>
        </div>
    );
}