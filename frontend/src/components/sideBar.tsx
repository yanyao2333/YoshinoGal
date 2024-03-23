import {AdjustmentsHorizontalIcon, BuildingLibraryIcon,} from '@heroicons/react/24/outline'

type SingleNavigation = {
    name: string;
    href: string;
    icon: any;
    current: boolean;
    count?: number;
}

const navigation: SingleNavigation[] = [
    {name: '游戏库', href: '/', icon: BuildingLibraryIcon, current: true},
    {name: '设置', href: '#', icon: AdjustmentsHorizontalIcon, current: false},
]

function classNames(...classes: string[]) {
    return classes.filter(Boolean).join(' ')
}

export default function SideBar() {
    return (
        <div className="flex grow flex-col gap-y-5 overflow-y-auto border-r border-gray-200 bg-white px-6">
            <div className="flex h-16 shrink-0 items-center">
                <img
                    className="h-8 w-auto"
                    src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSfvW3DwpzD8iXIC4TwjXGuFpaNgdHgXWGCkpa8Dh01yA&s"
                    alt="Your Company"
                />
            </div>
            <nav className="flex flex-1 flex-col">
                <ul role="list" className="flex flex-1 flex-col gap-y-7">
                    <li>
                        <ul role="list" className="-mx-2 space-y-1">
                            {navigation.map((item) => (
                                <li key={item.name}>
                                    <a
                                        href={item.href}
                                        className={classNames(
                                            item.current
                                                ? 'bg-gray-50 text-indigo-600'
                                                : 'text-gray-700 hover:text-indigo-600 hover:bg-gray-50',
                                            'group flex gap-x-3 rounded-md p-2 text-sm leading-6 font-semibold'
                                        )}
                                    >
                                        <item.icon
                                            className={classNames(
                                                item.current ? 'text-indigo-600' : 'text-gray-400 group-hover:text-indigo-600',
                                                'h-6 w-6 shrink-0'
                                            )}
                                            aria-hidden="true"
                                        />
                                        {item.name}
                                        {item.count ? (
                                            <span
                                                className="ml-auto w-9 min-w-max whitespace-nowrap rounded-full bg-white px-2.5 py-0.5 text-center text-xs font-medium leading-5 text-gray-600 ring-1 ring-inset ring-gray-200"
                                                aria-hidden="true"
                                            >
                        {item.count}
                      </span>
                                        ) : null}
                                    </a>
                                </li>
                            ))}
                        </ul>
                    </li>
                    <li className="-mx-6 mt-auto">
                        <a
                            className="flex items-center gap-x-4 px-6 py-3 text-sm font-semibold leading-6 text-gray-900"
                        >
                            <img
                                className="h-8 w-8 rounded-full bg-gray-50"
                                src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSfvW3DwpzD8iXIC4TwjXGuFpaNgdHgXWGCkpa8Dh01yA&s"
                                alt=""
                            />
                            <span className="select-none" aria-hidden="true">Yoshino Gal</span>
                        </a>
                    </li>
                </ul>
            </nav>
        </div>
    )
}
