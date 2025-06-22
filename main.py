from typing import Annotated
from typing_extensions import TypedDict

from langgraph.graph.message import add_messages  # type: ignore
from langgraph.graph import StateGraph, START  # type: ignore
from langchain_google_genai import ChatGoogleGenerativeAI  # type: ignore


class OrderState(TypedDict):
    # Tells LangGraph to add (append) messages instead of replacing them.
    messages: Annotated[list, add_messages]
    order: list[str]
    finished: bool


BARISTABOT_SYSINT = (
    "system",  # 'system' indicates the message is a system instruction.
    "You are a BaristaBot, an interactive cafe ordering system. A human will talk to you about the "
    "available products you have and you will answer any questions about menu items (and only about "
    "menu items - no off-topic discussion, but you can chat about the products and their history). "
    "The customer will place an order for 1 or more items from the menu, which you will structure "
    "and send to the ordering system after confirming the order with the human. "
    "\n\n"
    "Add items to the customer's order with add_to_order, and reset the order with clear_order. "
    "To see the contents of the order so far, call get_order (this is shown to you, not the user) "
    "Always confirm_order with the user (double-check) before calling place_order. Calling confirm_order will "
    "display the order items to the user and returns their response to seeing the list. Their response may contain modifications. "
    "Always verify and respond with drink and modifier names from the MENU before adding them to the order. "
    "If you are unsure a drink or modifier matches those on the MENU, ask a question to clarify or redirect. "
    "You only have the modifiers listed on the menu. "
    "Once the customer has finished ordering items, Call confirm_order to ensure it is correct then make "
    "any necessary updates and then call place_order. Once place_order has returned, thank the user and "
    "say goodbye!"
    "\n\n"
    "If any of the tools are unavailable, you can break the fourth wall and tell the user that "
    "they have not implemented them yet and should keep reading to do so.",
)

WELCOME_MSG = (
    "Welcome to the BaristaBot cafe. Type `q` to quit. How may I serve you today?"
)


llm = ChatGoogleGenerativeAI(model="gemini-2.5-flash")


def chatbot(state: OrderState) -> OrderState:
    message_history = [BARISTABOT_SYSINT] + state["messages"]
    return {"messages": [llm.invoke(message_history)]}


def build_graph():
    graph_builder = StateGraph(OrderState)
    graph_builder.add_node("chatbot", chatbot)
    graph_builder.add_edge(START, "chatbot")
    chat_graph = graph_builder.compile()
    return chat_graph


def main():
    chat_graph = build_graph()
    user_msg = "Hello, what can you do?"

    state = chat_graph.invoke({"messages": [user_msg]})
    for msg in state["messages"]:
        print(f"{type(msg).__name__}: {msg.content}")


if __name__ == "__main__":
    main()
