import React, { Component, Fragment } from 'react';
import { Button, Divider, Tooltip, Input, Modal, Table } from 'antd';

import FormModal from './Form';
import ToolBar from '../../components/ToolBar';
import { modulePaginate, moduleCreate, moduleUpdate } from '../../../services/appModule';

const Search = Input.Search;

class DataList extends Component {
    state = {
        data: {
            list: [],
            pagination: {},
        },
        loading: false,

        modalVisible: false,
        currentItem: {},
        modalType: 'create',
        search: '',
    }

    componentDidMount() {
        this.onQueryAll();
    }

    onQueryAll = (param = {}) => {
        this.setState({
            loading: true,
        });

        modulePaginate(param).then(response => {
            if(response.success){
                this.setState({
                    data: {
                        list: response.data,
                        pagination: {
                            current: response.current_page,
                            total: response.total,
                            pageSize: response.limit,
                        },
                    },
                    
                })
            } else {
                Modal.error({ title: 'Institución', content: response.message });
            }
        }).finally(e=>{
            this.setState({
                loading: false,
                modalVisible: false,
            });
        });
    }

    onShowModal = (modalType, currentItem = {}) => {
        this.setState({
            modalVisible: true,
            currentItem: currentItem,
            modalType: modalType,
        })
    };

    onSaveModal = values => {
        const modalType = this.state.modalType;
        this.setState({
            loading: true,
        });

        if(modalType === 'create'){
            moduleCreate(values).then(response => {
                if(response.success){
                    this.onQueryAll();
                }else{
                    Modal.error({ title: 'Institución', content: response.message });
                }
            }).finally(e=>{
                this.setState({
                    loading: false,
                });
            });
        }else if(modalType === 'update'){
            moduleUpdate({
                ...values,
                id: this.state.currentItem.id,
            }).then(response => {
                if(response.success){
                    this.onQueryAll();
                }else{
                    Modal.error({ title: 'Institución', content: response.message });
                }
            }).finally(e=>{
                this.setState({
                    loading: false,
                });
            });
        }
    }

    onCancelModal = () => {
        this.setState({
            modalVisible: false,
            currentItem: {},
            modalType: 'create',
            loading: false,
        })
    }

    handleTableChange = ({current,pageSize}) => {
        this.onQueryAll({
            current_page: current,
            limit: pageSize,
            search: this.state.search,
        });
    }

    onSearch = (e) => {
        this.setState({
            search: e.target.value,
        });
        this.onQueryAll({
            search: e.target.value,
        });
    }

    render() {
        const { loading } = this.state;

        const onDelete = (param) => {
            console.log(param);
        }

        const columns = [
            {
                title: 'Nombre',
                dataIndex: 'name',
                key: 'name',
            },
            {
                title: 'Denominación',
                dataIndex: 'description',
                key: 'description',
            },
            {
                title: 'Es nuevo',
                dataIndex: 'is_new',
                key: 'is_new',
            },
            {
                title: 'Accion',
                key: 'accion',
                width: '150px',
                render: (a, record) => {
                    return (
                        <div>
                            <Tooltip title="Editar">
                                <Button
                                    icon="edit"
                                    shape="circle"
                                    onClick={() => this.onShowModal('update', record)}
                                />
                            </Tooltip>
                            <Divider type="vertical"/>
                            <Tooltip title="Eliminar">
                                <Button
                                    icon="delete"
                                    shape="circle"
                                    onClick={() => {
                                        Modal.confirm({
                                            title: '¿Estás seguro de eliminar este registro?',
                                            content: a.name,
                                            okText: 'SI',
                                            okType: 'danger',
                                            cancelText: 'NO',
                                            onOk() {
                                                onDelete({ id: a.id });
                                            },
                                        });
                                    }}
                                />
                            </Tooltip>
                        </div>
                    );
                },
            },
        ];

        const paginationProps = {
            showSizeChanger: true,
            showQuickJumper: true,
            ...this.state.data.pagination,
        };

        return (
            <Fragment>
                <ToolBar
                    title="App módulos"
                    rightActions={[
                        <Tooltip title="Nuevo">
                            <Button
                                type="primary"
                                loading={loading}
                                onClick={() => this.onShowModal('create')}
                                icon="plus"
                            />
                        </Tooltip>,
                        <Divider type="vertical" />,
                        <Tooltip title="Regrescar">
                            <Button
                                icon="reload"
                                loading={loading}
                                onClick={()=>this.onQueryAll()}
                            />
                        </Tooltip>
                    ]}
                />
                <Search
                    style={{marginBottom: '16px'}}
                    placeholder="Buscar usuario"
                    value={this.state.search}
                    onChange={this.onSearch}
                />
                <Table
                    columns={columns}
                    size="small"
                    loading={loading}
                    rowKey={record => record.id}
                    dataSource={this.state.data.list} 
                    pagination={paginationProps} 
                    onChange={this.handleTableChange}
                />
                <FormModal 
                    loading={this.state.loading}
                    modalVisible={this.state.modalVisible}
                    currentItem={this.state.currentItem}
                    onSave={this.onSaveModal}
                    onCancel={this.onCancelModal}
                />
            </Fragment>
        );
    }
}

export default DataList;
