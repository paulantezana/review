import React from "react";
import { enquireScreen } from "enquire-js";
import { Link } from "gatsby";

import Header from "./Header";
import { Row, Col, Drawer, Icon, Menu } from "antd";

const TutorialMenu = ({tutorial = []}) => {
    return (
        <Menu>
            {
                tutorial.map((item,key) => (
                    <Menu.Item key={key}>
                        <Link to={item.node.fields.slug}>{item.node.frontmatter.title}</Link>
                    </Menu.Item>
                ))
            }
        </Menu>
    )
}

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
        const { children, tutorial } = this.props;

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
                                    <TutorialMenu tutorial={tutorial.edges} />
                                    <div className={`ApiMenuToggle ${!this.state.visibleMenu ? 'shadow' : ''}`} onClick={this.onToggle} >
                                        <Icon type={this.state.visibleMenu ? 'close' : 'menu-unfold'} />
                                    </div>
                                    
                                </Drawer>
                            ) : (
                                <TutorialMenu tutorial={tutorial.edges}/>
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
