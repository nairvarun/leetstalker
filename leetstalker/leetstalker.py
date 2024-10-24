import argparse
import asyncio
import json
import aiohttp
import tomllib
from sys import stderr
from pathlib import Path
from dataclasses import dataclass, asdict
from heapq import heapify, heappop
from rich.console import Console
from rich.table import Table
from rich.prompt import Confirm
from importlib.metadata import version, PackageNotFoundError


@dataclass
class Config:
    """
    loads configuration from toml file.

    default location: ~/.config/leetstalker/config.toml

    example config
    --------------
    users = [
        'nairvarun',
        'gultandon',
    ]
    """

    users: list[str]


@dataclass(frozen=True)
class Query:
    """holds graphql endpoint and query to fetch user data from leetcode"""

    url: str = "https://leetcode.com/graphql/"
    __query: str = """
        query userPublicProfileAndProblemsSolved($username: String!) {
        userContestRanking(username: $username) {
            rating
            globalRanking
            topPercentage
        }
        matchedUser(username: $username) {
            profile {
                ranking
            }
            submitStatsGlobal {
                acSubmissionNum {
                    difficulty
                    count
                }
            }
        }
    }"""

    @classmethod
    def get_query_data(cls, username):
        """substitutes variable and generate graphql query to fetch user data from leetcode"""
        return {"query": cls.__query, "variables": {"username": username}}


@dataclass(order=True)
class User:
    """holds user data fetched from leetcode"""

    ranking: int  # this will determine the order
    username: str

    # question metrics
    questions_solved: int
    hard_questions_solved: int
    medium_questions_solved: int
    easy_questions_solved: int

    # contest metrics
    contest_rating: float | None
    contest_ranking: int | None
    contest_percentile: float | None

    def __init__(self, username: str, data: dict, error: str | None = None):
        if error is not None:
            print(error, file=stderr)

        self.username = username
        if data == {}:
            return

        data = data["data"]
        if data["userContestRanking"] is not None:
            self.contest_rating = data["userContestRanking"]["rating"]
            self.contest_ranking = data["userContestRanking"]["globalRanking"]
            self.contest_percentile = data["userContestRanking"]["topPercentage"]
        else:
            self.contest_rating = None
            self.contest_ranking = None
            self.contest_percentile = None

        self.ranking = data["matchedUser"]["profile"]["ranking"]
        self.questions_solved = data["matchedUser"]["submitStatsGlobal"][
            "acSubmissionNum"
        ][0]["count"]
        self.easy_questions_solved = data["matchedUser"]["submitStatsGlobal"][
            "acSubmissionNum"
        ][1]["count"]
        self.medium_questions_solved = data["matchedUser"]["submitStatsGlobal"][
            "acSubmissionNum"
        ][2]["count"]
        self.hard_questions_solved = data["matchedUser"]["submitStatsGlobal"][
            "acSubmissionNum"
        ][3]["count"]


def get_configuration(config: Path) -> Config:
    if not config.exists():
        # offer to create config file
        resp = Confirm.ask(
            f"Create configuration file at {config}?",
            default="y",
        )
        if resp == "y":
            config.parent.mkdir(parents=True, exist_ok=True)
            with config.open(mode="w") as f:
                f.write(
                    """users = [
'larryNY',
'nairvarun',
]"""
                )
        else:
            exit(0)

    with config.open(mode="rb") as f:
        try:
            configuration: dict[str, list[str]] = tomllib.load(f)
        except Exception as e:
            print(e, file=stderr)
            exit(1)

    try:
        return Config(**configuration)
    except Exception as e:
        print("invalid config", file=stderr)
        exit(1)


async def get_data(session, username) -> User:
    try:
        async with session.post(
            Query.url, json=Query.get_query_data(username)
        ) as response:
            if response.status == 200:
                data = await response.json()
                if "errors" not in data:
                    return User(username, data, None)
                else:
                    return User("", {}, f"User {username} not found.")
            else:
                return User(
                    "",
                    {},
                    f"Failed to fetch data for {username}. Status code: {response.status}",
                )
    except Exception as e:
        return User("", {}, f"Failed to fetch data for {username}. {e}")


def print_table(responses: list[User], console: Console):
    heapify(responses)

    table = Table()

    table.add_column("Username", style="white")
    table.add_column("Ranking", style="magenta")
    table.add_column("Questions", style="blue")
    table.add_column("Hard", style="red")
    table.add_column("Medium", style="yellow")
    table.add_column("Easy", style="green")
    table.add_column("Contest Ranking", style="magenta")
    table.add_column("Contest Rating", style="blue")

    while responses:
        user: User = heappop(responses)
        table.add_row(
            user.username,
            str(user.ranking),
            str(user.questions_solved),
            str(user.hard_questions_solved),
            str(user.medium_questions_solved),
            str(user.easy_questions_solved),
            f"{user.contest_ranking:.02f}" if user.contest_ranking is not None else "-",
            f"{user.contest_rating:.02f}" if user.contest_rating is not None else "-",
        )

    console.print(table)


def print_json(responses: list[User]):
    responses = [asdict(user) for user in sorted(responses)]
    print(json.dumps(responses, indent=4))


async def leetstalker(users: list, config: Path, output: str) -> int:
    if not users:
        config: Config = get_configuration(config)
    else:
        config = Config(users)

    async with aiohttp.ClientSession() as session:
        with Console() as console:
            with console.status("Fetching data..."):
                tasks = (get_data(session, user) for user in config.users)

                responses = await asyncio.gather(*tasks)
                responses = [i for i in responses if i.username != ""]

                if output == "table":
                    print_table(responses, console)
                elif output == "json":
                    print_json(responses)


def root_handler(args):
    try:
        asyncio.run(leetstalker(args.users, args.config, args.output))
    except KeyboardInterrupt:
        exit(1)
    except Exception as e:
        print(e)
        exit(1)


def get_version():
    try:
        return version("leetstalker")
    except PackageNotFoundError:
        return "unknown"


def main():
    parser = argparse.ArgumentParser(
        prog="leetstalker",
        epilog="(https://github.com/nairvarun/leetstalker)",
        conflict_handler="resolve",
    )
    parser.add_argument(
        "--version", action="version", version=f"%(prog)s {get_version()}"
    )
    parser.add_argument("users", nargs="*")
    parser.add_argument(
        "--config",
        help="path to configuration file",
        type=Path,
        default=Path.joinpath(Path.home(), ".config", "leetstalker", "config.toml"),
    )
    parser.add_argument(
        "--output",
        help="specify output format",
        choices=["table", "json"],
        type=str,
        default="table",
    )
    parser.set_defaults(handler=root_handler)
    args = parser.parse_args()
    args.handler(args)


if __name__ == "__main__":
    main()
