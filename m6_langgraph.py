from typing import Annotated
from typing_extensions import TypedDict
from langgraph.graph import StateGraph, START, END
from langgraph.graph.message import add_messages
from langchain_anthropic import ChatAnthropic
from langchain_core.messages import SystemMessage, HumanMessage, ToolMessage
from lc_tools import LC_TOOLS, TOOL_BY_NAME

class State(TypedDict):
    messages: Annotated[list, add_messages]

SYSTEM = (
   "You are an instructor that teaches information from the personal knowledge vault (Cortex)."
   "Answer the user's question using ONLY what you find in the vault."
   "Start by reading 'index.md' (the catalog), then open the specific pages it points to."
   "Cite the file paths you drew from.  If the vault does not cover it, say so."
)

model = ChatAnthropic(model="claude-opus-4-8", max_tokens=16000).bind_tools(LC_TOOLS)
critic = ChatAnthropic(model="claude-opus-4-8", max_tokens=1024)

def agent_node(state: State) -> dict:
    ai = model.invoke([SYSTEM] + state["messages"])
    return {"messages": [ai]}

def tools_node(state: State) -> dict:
    ai = state["messages"][-1]
    results=[]
    for call in ai.tool_calls:
        print(f"   [tools] {call['name']}({call['args']})")
        output = TOOL_BY_NAME[call["name"]].invoke(call["args"])
        results.append(ToolMessage(content=output, tool_call_id=call["id"]))
    return {"messages": results}

def verify_node(state: State) -> dict:
    answer = state["messages"][-1].text
    # give the critic the evidence, not just the answer - the tool results are already sitting in state.  A critic that sees only the answer grades on vibes.
    sources = "\n\n".join(m.content for m in state["messages"] if isinstance(m, ToolMessage))
    verdict = critic.invoke([HumanMessage(
        "Below is an answer and the raw tool output it was built from. "
        "Does the answer cite at least one valut file path, and is every claim supported by the sources? "
        f"Reply PASS or FAIL with one sentence.\n\nAnswer: \n{answer}\n\nSources:\n{sources}"
    )])
    print(f"  [verify] {verdict.text}")
    return {}

def route(state: State) -> str:
    ai = state["messages"][-1]
    return "tools" if ai.tool_calls else "verify"

graph = StateGraph(State)
graph.add_node("agent", agent_node)
graph.add_node("tools", tools_node)
graph.add_node("verify", verify_node)
graph.add_edge(START, "agent")
graph.add_conditional_edges("agent", route, {"tools": "tools", "verify": "verify"})
graph.add_edge("tools", "agent") # after running tools, think again
graph.add_edge("verify", END)
app = graph.compile()
print(app.get_graph().draw_mermaid())

if __name__ == "__main__":
    result = app.invoke(
        {"messages": [HumanMessage("What does the vault say about when Not to use an agent?")]},
        config={"recursion_limit": 25},
    )
    print("\n=== ANSWER ===")
    print(result["messages"][-1].text)
