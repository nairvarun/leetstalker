import asyncio
import aiohttp

class Leetcode():
  def __init__(self):
    self.url: str = "https://leetcode.com/graphql/"
    self.query: str = """
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

  async def get_user_data(self, username: str):
    json_data: dict = {
      'query': self.query,
      'variables': {
        'username': username
      }
    }
    async with aiohttp.ClientSession() as session:
      async with session.post(self.url, json=json_data) as response:
        # handle errors
        # {'errors': [{'message': 'User matching query does not exist.', 'locations': [{'line': 3, 'column': 17}], 'path': ['userContestRanking']}, {'message': 'That user does not exist.', 'locations': [{'line': 8, 'column': 17}], 'path': ['matchedUser'], 'extensions': {'handled': True}}], 'data': {'userContestRanking': None, 'matchedUser': None}}
        if response.status == 200:
          data = await response.json()
          print(data)
        else:
          print(response.status)

class Configuration():
  ...

async def main() -> int:
  lc = Leetcode()
  await lc.get_user_data('gultandon')
  return 0

if __name__ == "__main__":
  exit(asyncio.run(main()))