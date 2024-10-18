import asyncio
import aiohttp
import tomllib
from pathlib import Path
from dataclasses import dataclass
from heapq import heapify, heappop
from rich.console import Console
from rich.table import Table


@dataclass
class Config:
  users: list[str]


@dataclass(frozen=True)
class Query:
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
    }
  """

  @classmethod
  def get_query_data(cls, username):
    return {
      'query': cls.__query,
      'variables': {
        'username': username
      }
    }


@dataclass(order=True)
class User:
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

  def __init__(self, username, data):
    self.username = username
    data = data['data']
    if data['userContestRanking'] is not None:
      self.contest_rating = data['userContestRanking']['rating']
      self.contest_ranking = data['userContestRanking']['globalRanking']
      self.contest_percentile = data['userContestRanking']['topPercentage']
    else:
      self.contest_rating = None
      self.contest_ranking = None
      self.contest_percentile = None

    # TODO: handle usr not exists err
    self.ranking = data['matchedUser']['profile']['ranking']
    self.questions_solved = data['matchedUser']['submitStatsGlobal']['acSubmissionNum'][0]['count']
    self.easy_questions_solved = data['matchedUser']['submitStatsGlobal']['acSubmissionNum'][1]['count']
    self.medium_questions_solved = data['matchedUser']['submitStatsGlobal']['acSubmissionNum'][2]['count']
    self.hard_questions_solved = data['matchedUser']['submitStatsGlobal']['acSubmissionNum'][3]['count']


def get_configuration() -> Config:
  configuration_file: Path = Path.joinpath(Path.home(), '.config', 'leetstalker', 'config.toml')
  with configuration_file.open(mode='rb') as f:
    configuration: dict[str, list[str]] = tomllib.load(f)
  return Config(**configuration)


async def get_data(session, username) -> User:
    async with session.post(Query.url, json=Query.get_query_data(username)) as response:
      # handle errors
        # {'errors': [{'message': 'User matching query does not exist.', 'locations': [{'line': 3, 'column': 17}], 'path': ['userContestRanking']}, {'message': 'That user does not exist.', 'locations': [{'line': 8, 'column': 17}], 'path': ['matchedUser'], 'extensions': {'handled': True}}], 'data': {'userContestRanking': None, 'matchedUser': None}}
      if response.status == 200:
        data = await response.json()
        if 'errors' not in data:
          return User(username, data)
        else:
          # TODO: handle gracefully
          print(data['errors'])
      else:
        # TODO: handle gracefully
        print(response.status)

async def main() -> int:
  config: Config = get_configuration()
  async with aiohttp.ClientSession() as session:
    tasks = [get_data(session, user) for user in config.users]

    responses = await asyncio.gather(*tasks)
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
        str(user.contest_ranking) if user.contest_ranking is not None else '-',
        str(user.contest_rating) if user.contest_rating is not None else '-',
      )

    console = Console()
    console.print(table)

  return 0

if __name__ == "__main__":
  exit(asyncio.run(main()))