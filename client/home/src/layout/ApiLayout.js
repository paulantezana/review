import React from "react"
import { enquireScreen } from "enquire-js"

import Header from "./Header"
import ApiMenu from "./ApiMenu"

import { Row, Col, Drawer, Icon } from "antd"

// let isMb

// enquireScreen(b => {
//     isMb = b
// })

class SiteLayout extends React.Component {

    constructor(props){
        super(props)
        this.state = {
            isMobile : false,
            visibleMenu : true,
        }
    }

    componentDidMount() {
        enquireScreen(b => {
            this.setState({
                isMobile: !!b,
            })
        })
    }

    onClose = () => {
        this.setState({
            visibleMenu: false,
        })
    }

    onToggle = () => {
        this.setState({
            visibleMenu: !this.state.visibleMenu,
        })
    }

    render() {
        const { children } = this.props
        console.log(this.state)
        return (
            <div>
                <Header isMobile={this.state.isMobile} />
                <Row style={{ marginTop: "32px" }}>
                    <Col lg={6} xl={5} xxl={4}>
                        {
                            this.state.isMobile ? (
                                <Drawer
                                    closable={false}
                                    placement="left"
                                    bodyStyle={{padding: 0}}
                                    onClose={this.onClose}
                                    visible={this.state.visibleMenu}
                                >
                                    <ApiMenu />
                                    <div className={`ApiMenuToggle ${!this.state.visibleMenu ? 'shadow' : ''}`} onClick={this.onToggle} >
                                        <Icon type={this.state.visibleMenu ? 'close' : 'menu-unfold'} />
                                    </div>
                                    
                                </Drawer>
                            ) : (
                                <ApiMenu />
                            )
                        }
                    </Col>
                    <Col lg={18} xl={19} xxl={20}>
                        {children}
                    </Col>
                </Row>
            </div>
        )
    }
}

export default SiteLayout
