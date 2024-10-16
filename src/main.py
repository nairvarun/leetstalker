import asyncio
import aiohttp
import tomllib
from pathlib import Path
from dataclasses import dataclass



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


def get_configuration() -> Config:
  configuration_file: Path = Path.joinpath(Path.home(), '.config', 'leetstalker', 'config.toml')
  with configuration_file.open(mode='rb') as f:
    configuration: dict[str, list[str]] = tomllib.load(f)
  return Config(**configuration)


async def get_data(session, user):
    async with session.post(Query.url, json=Query.get_query_data(user)) as response:
      # handle errors
        # {'errors': [{'message': 'User matching query does not exist.', 'locations': [{'line': 3, 'column': 17}], 'path': ['userContestRanking']}, {'message': 'That user does not exist.', 'locations': [{'line': 8, 'column': 17}], 'path': ['matchedUser'], 'extensions': {'handled': True}}], 'data': {'userContestRanking': None, 'matchedUser': None}}
      if response.status == 200:
        data = await response.json()
        return data
      else:
        print(response.status)

async def main() -> int:
  config: Config = get_configuration()
  async with aiohttp.ClientSession() as session:
    tasks = (get_data(session, user) for user in config.users)

    responses = await asyncio.gather(*tasks)
    for response in responses:
      print(response)
  return 0

if __name__ == "__main__":
  exit(asyncio.run(main()))