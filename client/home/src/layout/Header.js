import React from "react";
import { Row, Col, Icon, Menu, Popover, Avatar, Dropdown } from "antd";
import { Link, navigate } from "gatsby";

import { enquireScreen } from "enquire-js";
import { getToken, getAuthorityLicense, destroy } from '../utils/authority';
import { service } from '../utils/config';

import logoSvg from '../images/logo.svg';

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

    onMenuAvatar = ({key}) => {
        switch (key) {
            case 'profile':
                navigate('/admin/profile');
                break;
            case 'logout':
                destroy();
                navigate('/admin');
                break;
            default:
                break;
        }
    }

    render() {
        const { menuMode, menuVisible } = this.state
        const { user } = getAuthorityLicense();

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
                {
                    getToken() && (
                        <Menu.SubMenu
                            key="subAdmin"
                            title="Admin"
                        >
                            <Menu.Item key="admin">
                                <Link to="/admin">Deshboard</Link>
                            </Menu.Item>
                            <Menu.Item key="module">
                                <Link to="/admin/module">Módulos</Link>
                            </Menu.Item>
                            <Menu.Item key="admin">
                                <Link to="/admin/function">Funciones</Link>
                            </Menu.Item>
                            <Menu.Item key="institution">
                                <Link to="/admin/institution">Instituciones</Link>
                            </Menu.Item>
                        </Menu.SubMenu>
                    )
                }
                {menuMode === "inline" && (
                    <Menu.Item key="preview">
                        {/* <a target="_blank" href="http://preview.pro.ant.design/" rel="noopener noreferrer">
                预览
              </a> */}
                    </Menu.Item>
                )}
            </Menu>
        )

        const menuAvatar = (
            <Menu selectedKeys={[]} onClick={this.onMenuAvatar}>
              <Menu.Item key="profile">
                <Icon type="user" />
                Perfil
              </Menu.Item>
              <Menu.Divider />
              <Menu.Item key="logout">
                <Icon type="logout" />
                Serrar sesión
              </Menu.Item>
            </Menu>
        );


        return (
            <div className="Header color">
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
                            <img src={logoSvg} alt="logo" />
                            <span>Documentacion</span>
                        </Link>
                    </Col>
                    <Col xxl={20} xl={19} lg={16} md={16} sm={0} xs={0}>
                        <div className="header-meta">
                            <div id="preview">
                                {
                                    getToken() ? (
                                        <Dropdown overlay={menuAvatar}>
                                            <Avatar
                                                src={`${service.path}/${user.avatar}`}
                                                alt="avatar"
                                            />
                                        </Dropdown>
                                    ) : (
                                        <Avatar icon="user" style={{ backgroundColor: '#FAAD14' }} onClick={()=>navigate('/admin/profile')}/>
                                    )
                                }
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
