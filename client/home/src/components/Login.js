import React, { Component } from 'react';
import { navigate } from "gatsby";
import { Form, Icon, Input, Button, Checkbox, Modal } from 'antd';

import { setAuthority, getToken } from "../utils/authority";
import { login } from '../services/appUser';

class LoginForm extends Component {
    constructor(props) {
        super(props);
        this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleSubmit(e) {
        e.preventDefault();
        this.props.form.validateFields((err, values) => {
            if (!err) {
                login({ ...values }).then(response=>{
                    if(response.success){
                        setAuthority({
                            token: response.data.token,
                            remember: values.remember,
                            user: response.data.app_user,
                        });
                        navigate('/admin/profile');
                    }else{
                        Modal.error({ title: 'Login', content: response.message });
                    }
                });
            }
        });
    }

    render() {
        const { getFieldDecorator } = this.props.form;
        const { loading } = this.props;
        if (getToken()) {
            navigate(`/admin/profile`)
        }

        return (
            <div style={{maxWidth: '300px', margin: '5rem auto'}}>
                <Form onSubmit={this.handleSubmit}>
                    <Form.Item>
                        {getFieldDecorator('user_name', {
                            rules: [ { required: true,  message: 'Por favor ingrese su nombre de usuario!' } ],
                        })(
                            <Input
                                prefix={<Icon type="user" style={{ color: 'rgba(0,0,0,.25)' }} />}
                                placeholder='Usuario'
                            />
                        )}
                    </Form.Item>
                    <Form.Item>
                        {getFieldDecorator('password', {
                            rules: [ { required: true, message: 'Por favor ingrese su contraseña!' } ],
                        })(
                            <Input.Password
                                prefix={<Icon type="lock" style={{ color: 'rgba(0,0,0,.25)' }} />}
                                placeholder="Contraseña"
                            />
                        )}
                    </Form.Item>
                    <Form.Item>
                        {getFieldDecorator('remember', {
                            valuePropName: 'checked',
                            initialValue: false,
                        })(
                            <Checkbox>
                                Recuérdame
                            </Checkbox>
                        )}
                        <Button type="primary" loading={loading} htmlType="submit" block>
                            Iniciar sesión
                        </Button>
                    </Form.Item>
                </Form>
            </div>
        );
    }
}

const LoginPage = Form.create()(LoginForm);

export default LoginPage;