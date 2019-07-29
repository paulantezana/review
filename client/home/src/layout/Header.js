import React from "react"
import { Row, Col, Icon, Menu, Button, Popover } from "antd"
import { Link } from "gatsby"

import { enquireScreen } from "enquire-js"
import { DataIntances } from '../data/data';

const LOGO_URL = "https://assets.paulantezana.com/logo-circle-white-purple.svg"

class Header extends React.Component {
    state = {
        menuVisible: false,
        menuMode: "horizontal",
    }

    componentDidMount() {
        enquireScreen(b => {
            this.setState({ menuMode: b ? "inline" : "horizontal" })
        })
    }

    toggleMenu = ()=>{
        this.setState({
            menuVisible: !this.state.visible
        })
    }

    render() {
        const { menuMode, menuVisible } = this.state

        const menu = (
            <Menu mode={menuMode} id="nav" key="nav">
                <Menu.Item key="home">
                    <Link to="/">Inicio</Link>
                </Menu.Item>
                <Menu.Item key="tutorial">
                    <Link to="/tutorial">Tutorial</Link>
                </Menu.Item>
                <Menu.Item key="documentacion">
                    <Link to="/documentacion">Documentación</Link>
                </Menu.Item>
                {menuMode === "inline" && (
                    <Menu.Item key="preview">
                        {/* <a target="_blank" href="http://preview.pro.ant.design/" rel="noopener noreferrer">
                预览
              </a> */}
                    </Menu.Item>
                )}
            </Menu>
        )

        return (
            <div className="Header">
                {menuMode === "inline" ? (
                    <Popover
                        overlayClassName="popover-menu"
                        placement="bottomRight"
                        content={menu}
                        trigger="click"
                        visible={menuVisible}
                        arrowPointAtCenter
                        // onVisibleChange={this.onMenuVisibleChange}
                    >
                        <Icon
                            className="nav-phone-icon"
                            type="menu"
                            onClick={this.toggleMenu}
                        />
                    </Popover>
                ) : null}
                <Row className="Container">
                    <Col xxl={4} xl={5} lg={8} md={8} sm={24} xs={24}>
                        <Link to="/" className="Header-logo">
                            <img src={LOGO_URL} alt="logo" />
                            <span>Documentacion</span>
                        </Link>
                    </Col>
                    <Col xxl={20} xl={19} lg={16} md={16} sm={0} xs={0}>
                        <div className="header-meta">
                            <div id="preview">
                                <a
                                    href={DataIntances.find(item=>item.key === 3).url}
                                    target="_blanck"
                                >
                                    <Button icon="eye-o" type="primary">
                                        Ingresar
                                    </Button>
                                </a>
                            </div>
                            {menuMode === "horizontal" ? (
                                <div className="Menu">{menu}</div>
                            ) : null}
                        </div>
                    </Col>
                </Row>
            </div>
        )
    }
}

export default Header
