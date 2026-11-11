from langchain.agents import create_agent
from lc_tools import LC_TOOLS

SYSTEM = (
   "You are a research assistant for a personal knowledge vault (Cortex). "
   "Answer the user's question using ONLY what you find in the vault"
   "Start by reading 'index.md' (the catalog), then open the specific pages it points to.  "
   "Cite the file paths you drew from.  If the vault does not cover it, say so."
)

agent = create_agent(
   model="anthropic:claude-opus-4-8",
   tools=LC_TOOLS,
   system_prompt=SYSTEM,
)

result = agent.invoke({"messages": [{"role": "user", "content": "What does the vault say about cooking."}]})
print(result["messages"][-1].text)
