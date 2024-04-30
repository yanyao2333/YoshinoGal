import SideBar from "../components/sideBar";
import GamePosterWall from "../components/gamePosterWall";
import {GetPosterWall} from "../../wailsjs/go/app/Library";
import {useEffect, useState} from "react";
import ManualScrapeDialog from "../components/dialog";

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
            imageSrc: "data:image/jpg;base64," + result[game].poster_b64,
            id: result[game].game_id,
            href: '/game/' + result[game].game_id,
            imageAlt: result[game].game_name,
        });
    }

    return gamesShowList
}


export default function Home() {
    const [games, setGames] = useState<Game[]>([]);
    const [dialog, setDialog] = useState(false)

    useEffect(() => {
        getPosterWall().then(result => {
            setGames(result);
        }).catch(error => {
            console.error("Error fetching poster wall:", error);
        });
    }, []);
    if (dialog) {
        return (
            <div className="flex flex-row bg-gray-100">
                <ManualScrapeDialog/>
                <div className="flex flex-grow gap-y-5 w-72 h-screen overflow-y-auto">
                    <SideBar/>
                </div>

                <div className="flex-col overflow-y-auto h-screen bg-white divide-y">
                    <div className="flex pt-3 pb-3 gap-3 place-content-between">
                        <p className="self-center text-2xl text-gray-900 pl-3 font-black">
                            游戏库
                        </p>
                        <div className="pr-3">
                            <button
                                onClick={() => setDialog(false)}
                                type="button"
                                className="rounded-md bg-indigo-600 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
                            >
                                手动刮削
                            </button>
                        </div>
                    </div>
                    <GamePosterWall game={games}/>
                </div>
            </div>
        );
    }
    return (
        <div className="flex flex-row bg-gray-100">
            <div className="flex flex-grow gap-y-5 w-72 h-screen overflow-y-auto">
                <SideBar/>
            </div>

            <div className="flex-col overflow-y-auto h-screen bg-white divide-y">
                <div className="flex pt-3 pb-3 gap-3 place-content-between">
                    <p className="self-center text-2xl text-gray-900 pl-3 font-black">
                        游戏库
                    </p>
                    <div className="pr-3">
                        <button
                            onClick={() => setDialog(true)}
                            type="button"
                            className="rounded-md bg-indigo-600 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
                        >
                            手动刮削
                        </button>
                    </div>
                </div>
                <GamePosterWall game={games}/>
            </div>
        </div>
    );
}