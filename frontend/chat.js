export class chat extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            input: "",
        }
    }

    render() {
        let messages = [];
        this.props.messages.forEach((message, i) => {
            messages.push(React.createElement(
                "div", {
                    key: "message_"+i,
                    style: {
                        display: "flex",
                        flexDirection: "row",
                        justifyContent: "center",
                        alignItems: "center",
                        border: "1px solid black",
                        borderRadius: "5px",
                    },
                },
                message.name + ": " + message.text + " (" + new Date(message.sentAt).toLocaleTimeString() + " " + new Date(message.sentAt).toLocaleDateString() + ")"
            ))
        })
        messages.reverse();
        return React.createElement(
            "div", {
                id: "messages",
                onContextMenu: (e) => {
                    e.preventDefault();
                },
                style: {
                    fontFamily: "sans-serif",
                    borderRadius: "5px",
                    overflow: "auto",
                    margin: "auto",
                },
            },
            React.createElement(
                "div", {    
                    style: {
                        fontFamily: "sans-serif",
                        display: "flex",
                        flexDirection: "column",
                        borderRadius: "5px",
                        overflow: "auto",
                        height: "500px",
                        border: "1px solid black",
                    },
                },
                messages
            ),
            React.createElement(
                "input", {
                    id: "message",
                    type: "text",
                    onKeyPress: (e) => {
                        if (e.key === "Enter") {
                            this.props.WS_CONNECTION.send(this.props.constructMessage("addMessage", this.state.input));
                            this.setState({
                                input: "",
                            });
                        }
                    },
                    value: this.state.input,
                    onChange: (e) => {
                        this.setState({
                            input: e.target.value,
                        });
                    },
                    style: {
                        fontFamily: "sans-serif",
                        display: "flex",
                        flexDirection: "column" ,
                        border: "1px solid black",
                        borderRadius: "5px",
                        margin:"auto",
                        padding: "5px",
                    },
                }
            ),
        );
    }
}