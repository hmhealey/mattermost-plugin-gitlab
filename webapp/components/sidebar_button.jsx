import PropTypes from 'prop-types';
import React from 'react';
import {Tooltip, OverlayTrigger} from 'react-bootstrap';

export default class SidebarButton extends React.PureComponent {
	static propTypes = {
		count: PropTypes.number,
		href: PropTypes.string,
		icon: PropTypes.node.isRequired,
		onClick: PropTypes.func,
		style: PropTypes.object,
		tooltipId: PropTypes.string,
		tooltipPlacement: PropTypes.string,
		tooltipText: PropTypes.string,
	};

	render() {
		let href = this.props.href;
		if (href == null && this.props.onClick) {
			href = '#';
		}

		let buttonText;
		if (this.props.count != null) {
			buttonText = ' ' + this.props.count;
		}

		let button = (
			<a
                href={href}
                target='_blank'
                style={this.props.style}
            >
            	{this.props.icon}
                {buttonText}
            </a>
        );

        if (this.props.tooltipText) {
        	button = (
                <OverlayTrigger
	                id={this.props.tooltipId}
                    placement={this.props.tooltipPlacement}
                    overlay={<Tooltip>{this.props.tooltipText}</Tooltip>}
                >
                	{button}
                </OverlayTrigger>
            );
        }

        return button;
	}
}