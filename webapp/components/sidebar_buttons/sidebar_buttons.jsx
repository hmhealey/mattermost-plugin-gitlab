import PropTypes from 'prop-types';
import React from 'react';

import {makeStyleFromTheme, changeOpacity} from 'mattermost-redux/utils/theme_utils';

import SidebarButton from '../sidebar_button.jsx';

export default class SidebarButtons extends React.PureComponent {
    static propTypes = {
        theme: PropTypes.object.isRequired,
        gitlabURL: PropTypes.string.isRequired,
        connected: PropTypes.bool,
        gitlabUsername: PropTypes.string,
        gitlabUserId: PropTypes.number,
        assignedIssues: PropTypes.arrayOf(PropTypes.object),
        assignedMergeRequests: PropTypes.arrayOf(PropTypes.object),
        createdMergeRequests: PropTypes.arrayOf(PropTypes.object),
        todos: PropTypes.arrayOf(PropTypes.object),
        isTeamSidebar: PropTypes.bool,
        actions: PropTypes.shape({
            getAssignedIssues: PropTypes.func.isRequired,
            getAssignedMergeRequests: PropTypes.func.isRequired,
            getCreatedMergeRequests: PropTypes.func.isRequired,
            getTodos: PropTypes.func.isRequired,
        }).isRequired
    };

    static defaultProps = {
        assignedIssues: [],
        assignedMergeRequests: [],
        createdMergeRequests: [],
    };

    constructor(props) {
        super(props);

        this.state = {
            refreshing: false,
        };
    }

    componentDidMount() {
        if (this.props.connected) {
            this.getData();
        }
    }

    componentDidUpdate(prevProps) {
        if (this.props.connected && !prevProps.connected) {
            this.getData();
        }
    }

    getData = async (e) => {
        if (this.state.refreshing) {
            return;
        }

        if (e) {
            e.preventDefault();
        }

        this.setState({refreshing: true});
        await Promise.all([
            this.props.actions.getAssignedIssues(),
            this.props.actions.getAssignedMergeRequests(),
            this.props.actions.getCreatedMergeRequests(),
            this.props.actions.getTodos(),
        ]);
        this.setState({refreshing: false});
    }

    openConnectWindow = (e) => {
        e.preventDefault();
        window.open('/plugins/gitlab/oauth/connect', 'Connect Mattermost to GitLab', 'height=570,width=520');
    }

    render() {
        const style = getStyle(this.props.theme);
        const isTeamSidebar = this.props.isTeamSidebar;

        let container = style.containerHeader;
        let button = style.buttonHeader;
        let placement = 'bottom';
        if (isTeamSidebar) {
            placement = 'right';
            button = style.buttonTeam;
            container = style.containerTeam;
        }

        if (!this.props.connected) {
            if (isTeamSidebar) {
                return (
                    <SidebarButton
                        icon={<i className='fa fa-gitlab fa-2x'/>}
                        href='/plugins/gitlab/oauth/connect'
                        onClick={this.openConnectWindow}
                        style={button}
                        tooltipId='connectToGitLabTooltip'
                        tooltipPlacement={placement}
                        tooltipText='Connect to your GitLab account'
                    />
                )
            } else {
                return null;
            }
        }

        const refreshClass = this.state.refreshing ? ' fa-spin' : '';

        return (
            <div style={container}>
                <SidebarButton
                    href={this.props.gitlabURL + '/profile/applications'}
                    icon={<i className='fa fa-gitlab'/>}
                    style={button}
                />
                <SidebarButton
                    count={this.props.createdMergeRequests.length}
                    href={this.props.gitlabURL + '/dashboard/merge_requests?author_id=' + this.props.gitlabUserId}
                    icon={<i className='fa fa-code-fork'/>}
                    style={button}
                    tooltipId='createdMergeRequestsTooltip'
                    tooltipPlacement={placement}
                    tooltipText='Pull requests you created'
                />
                <SidebarButton
                    count={this.props.assignedMergeRequests.length}
                    href={this.props.gitlabURL + '/dashboard/merge_requests?assignee_id=' + this.props.gitlabUserId}
                    icon={<i className='fa fa-envelope'/>}
                    style={button}
                    tooltipId='assignedMergeRequestsTooltip'
                    tooltipPlacement={placement}
                    tooltipText="Pull requests you're assigned to"
                />
                <SidebarButton
                    count={this.props.assignedIssues.length}
                    href={this.props.gitlabURL + '/dashboard/issues?assignee_id=' + this.props.gitlabUserId}
                    icon={<i className='fa fa-sticky-note'/>}
                    style={button}
                    tooltipId='assignedIssuesTooltip'
                    tooltipPlacement={placement}
                    tooltipText="Issues you're assigned to"
                />
                <SidebarButton
                    count={this.props.todos.length}
                    href={this.props.gitlabURL + '/dashboard/todos'}
                    icon={<i className='fa fa-check-square'/>}
                    style={button}
                    tooltipId='todosTooltip'
                    tooltipPlacement={placement}
                    tooltipText="Your pending TODOs"
                />
                {/*<SidebarButton
                    count={this.props.unreads.length}
                    href={this.props.gitlabURL + '/pulls?q=is%3Aopen+mentions%3A' + this.props.gitlabUsername + '+archived%3Afalse'}
                    icon={<i className='fa fa-envelope'/>}
                    style={button}
                    tooltipId='unreadMessagesTooltip'
                    tooltipPlacement={placement}
                    tooltipText='Unread messages'
                />*/}
                <SidebarButton
                    icon={<i className={'fa fa-refresh' + refreshClass}/>}
                    onClick={this.getData}
                    style={button}
                    tooltipId='refreshTooltip'
                    tooltipPlacement={placement}
                    tooltipText='Refresh'
                />
            </div>
        );
    }
}

const getStyle = makeStyleFromTheme((theme) => {
    return {
        buttonTeam: {
            color: changeOpacity(theme.sidebarText, 0.6),
            display: 'block',
            marginBottom: '10px',
            width: '100%',
        },
        buttonHeader: {
            color: changeOpacity(theme.sidebarText, 0.6),
            flex: 1,
            textAlign: 'center',
            cursor: 'pointer',
        },
        containerHeader: {
            marginTop: '10px',
            marginBottom: '5px',
            display: 'flex',
            alignItems: 'center',
        },
        containerTeam: {
        },
    };
});
