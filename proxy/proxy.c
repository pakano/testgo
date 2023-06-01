#include<netdb.h>
#include<stdio.h>
#include<stdlib.h>
#include<unistd.h>
#include<string.h>
#include<errno.h>
#include<sys/socket.h>
#include<sys/types.h>
#include <sys/wait.h>
#include<signal.h>


#define bool int
#define TRUE 1
#define FALSE 0
#define OK 1
#define BUF_SIZE 1024

int remote_port, local_port;
int listen_fd;
char *remote_host;
struct sockaddr_in server_addr;

int create_server(int local_port);
int server_loop();
int handler_client(int connfd);
int create_conn_server(char* remote_host, int remote_port);
int redirect_data(int srcfd, int dstfd);
void sigchld_handler(int signal);

int main(int argc, char *argv[]) {
	int opt;
	char **pptr;
	bool show_version_flag = FALSE;
	bool remote_host_flag = FALSE;
	bool local_port_flag = FALSE;
	bool remote_port_flag = FALSE;
	char str[100];
	struct hostent *host_s;

	/*
	 * v For version
	 * l For localhost port
	 * h For remote ip
	 * p For remote port
	 */

	while ((opt = getopt(argc, argv, "vl:h:p:")) != -1) {
		switch (opt) {
		case 'v':
			show_version_flag = TRUE;
			break;
		case 'h':
			remote_host_flag = TRUE;
			remote_host = optarg;
			break;
		case 'l':
			local_port_flag = TRUE;
			local_port = atoi(optarg);
			break;
		case 'p':
			remote_port_flag = TRUE;
			remote_port = atoi(optarg);
			break;
		default:
			break;
		}
	}

	if (show_version_flag) {
		printf("TcpProxy Version 1.1\n");
	}

	if (!(local_port_flag && remote_port_flag && remote_port_flag)) {
		printf("Usage: tcpproxy -l 8080 -h 127.0.0.1 -p 80\n");
		exit(0);
	} else {
		printf("tcpproxy -l %d -h %s -p %d\n", local_port, remote_host,
				remote_port);
	}
	signal(SIGCHLD, sigchld_handler); // prevent ended children from becoming zombies
	if (create_server(local_port) == -1) {
		printf("create server() error.Exiting...");
		exit(-1);
	}
	server_loop();
	return 0;
}

void sigchld_handler(int signal) {
	while (waitpid(-1, NULL, WNOHANG) > 0)
		;
}

int create_server(int local_port) {
	listen_fd = socket(AF_INET, SOCK_STREAM, 0);
	if (listen_fd == -1) {
		printf("socket() error %d : %s\n", errno, strerror(errno));
		return -1;
	}
	server_addr.sin_port = htons(local_port);
	server_addr.sin_addr.s_addr = htonl(INADDR_ANY );
	if ((bind(listen_fd, (struct sockaddr *) &server_addr, sizeof(server_addr)))
			== -1) {
		printf("bind() error %d : %s \n", errno, strerror(errno));
		return -1;
	}
	if ((listen(listen_fd, 20)) == -1) {
		printf("listen() error %d : %s \n", errno, strerror(errno));
		return -1;
	}
	return listen_fd;
}

int server_loop() {
	int connfd;
	while (1) {
		connfd = accept(listen_fd, (struct sockaddr *) NULL, NULL );
		if (connfd == -1) {
			printf("accept() error %d : %s\n", errno, strerror(errno));
			return -1;
		}
		if (fork() == 0) {
			close(listen_fd);
			handler_client(connfd);
			exit(0);
		}
		close(connfd);
	}
	return OK;
}

int handler_client(int connfd) {
	int client_fd;
	client_fd = create_conn_server(remote_host, remote_port);
	if (fork() == 0) {
		redirect_data(connfd, client_fd);
		exit(0);
	}
	if (fork() == 0) {
		redirect_data(client_fd, connfd);
		exit(0);
	}
	close(connfd);
	close(client_fd);
	return OK;
}

int create_conn_server(char* remote_addr, int remote_port) {
	int client_fd;
	struct sockaddr_in client_addr;
	client_fd = socket(AF_INET, SOCK_STREAM, 0);
	if (client_fd == -1) {
		printf("socket errno %d : %s \n", errno, strerror(errno));
		return -1;
	}
	client_addr.sin_family = AF_INET;
	client_addr.sin_port = htons(remote_port);
	inet_pton(AF_INET, remote_addr, &client_addr.sin_addr);
	if ((connect(client_fd, (struct sockaddr *) &client_addr,
			sizeof(client_addr))) < 0) {
		printf("connect error %d : %s \n", errno, strerror(errno));
		return -1;
	}
	return client_fd;
}

int redirect_data(int srcfd, int dstfd) {
	char buf[BUF_SIZE];
	int n;
	while ((n = recv(srcfd, buf, BUF_SIZE, 0)) > 0) {
		send(dstfd, buf, n, 0);
	}
	close(srcfd);
	close(dstfd);
	return OK;
}