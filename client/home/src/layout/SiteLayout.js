import React from "react"
import { enquireScreen } from "enquire-js"

import Header from "./Header"
import "../scss/app.scss"

let isMobile

enquireScreen(b => {
    isMobile = b
})

class SiteLayout extends React.PureComponent {
    state = {
        isMobile,
    }

    componentDidMount() {
        enquireScreen(b => {
            this.setState({
                isMobile: !!b,
            })
        })
    }

    render() {
        const { children } = this.props
        return (
            <div>
                <Header isMobile={this.state.isMobile} />
                {children}
            </div>
        )
    }
}

export default SiteLayout
