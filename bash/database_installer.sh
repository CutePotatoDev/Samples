#!/bin/bash
trap finish 2

# Some colors and stuff.
RED="\e[31m"
YELLOW="\e[33m"
GREEN="\e[32m"
CYAN="\e[36m"
ORANGE="\e[38;5;208m"
ENDC="\e[0m"
UNDERLINE="\e[4m"
NORMAL="\e[24m"

CONF="$HOME/.myconf.cnf"
BACKUPS="./backups/"


function configure() {
    echo -e $CYAN
    echo "############################################################"
    echo "# You entered script configuration area, no change will be #"
    echo "# performed in your DB's.                                  #"
    echo "# Just some questions be asked about your host and DB's.   #"
    echo -e "############################################################"$ENDC


    #Login Server.
    echo -ne $CYAN"\nMySQL Login Server hostname (default localhost): "$ENDC
    read LS_DBHOST

    if [ -z "$LS_DBHOST" ]; then
        LS_DBHOST="localhost"
    fi

    echo -ne $CYAN"MySQL Login Server database name (default login): "$ENDC
    read LS_DB

    if [ -z "$LS_DB" ]; then
        LS_DB="login"
    fi

    echo -ne $CYAN"MySQL Login Server user (default root): "$ENDC
    read LS_USER

    if [ -z "$LS_USER" ]; then
        LS_USER="root"
    fi

    echo -ne $CYAN"Now script ask you enter passwor. Press any key to continue..."$ENDC
    read KEY

    mysql_config_editor set --login-path=$LS_DB --host=$LS_DBHOST --user=$LS_USER --password --skip-warn

    #Game Server
    echo -ne $CYAN"\nMySQL Game Server hostname (default $LS_DBHOST): "$ENDC
    read GS_DBHOST

    if [ -z "$GS_DBHOST" ]; then
        GS_DBHOST="localhost"
    fi

    echo -ne $CYAN"MySQL Game Server database name (default game): "$ENDC
    read GS_DB

    if [ -z "$GS_DB" ]; then
        GS_DB="game"
    fi

    echo -ne $CYAN"MySQL Game Server user (default $LS_USER): "$ENDC
    read GS_USER

    if [ -z "$GS_USER" ]; then
        GS_USER="root"
    fi

    echo -ne $CYAN"Now script ask you enter passwor. Press any key to continue..."$ENDC
    read KEY

    mysql_config_editor set --login-path=$GS_DB --host=$LS_DBHOST --user=$LS_USER --password --skip-warn

    saveConfig $1
}

# Mysql config already saved. Here just save default databases.
function saveConfig() {
    # Create empty file or drop if exist.
    echo -n > $CONF

# Megh not like indention.
cat <<EOF > $CONF
#Configuration settings for database installer script.
LS_DB=$LS_DB
GS_DB=$GS_DB
EOF

    chmod 600 $CONF

    echo ""
    echo -e $YELLOW"Configuration saved."$ENDC
}


function checkConfig() {
    FLAG=False

    if [ -e "$CONF" ] && [ -f "$CONF" ]; then
        # Read file.
        . $CONF
    else
        echo -e $YELLOW"Default database names credentials doesn't exist."$ENDC
        FLAG=True
    fi


    OUT=$( mysql_config_editor print --login-path=$LS_DB )

    if [ "$OUT" == "" ]; then
        echo -e $YELLOW"Login database credentials doesn't exist."$ENDC
        FLAG=True
    fi

    OUT=$( mysql_config_editor print --login-path=$GS_DB )

    if [ "$OUT" == "" ]; then
        echo -e $YELLOW"Game database credentials doesn't exist."$ENDC
        FLAG=True
    fi

    if [ $FLAG == True ]; then
        echo ""
        echo -e $RED"Data bases login credentials file, not found."$ENDC
        echo ""

        echo -ne $CYAN"You must create new credentials file, continue ??? (Y/n): "
        read CONT

        if [ "$CONT" == "Y" -o "$CONT" == "y" ]; then
            configure
        elif [ "$CONT" == "n" -o "$CONT" == "N" ]; then 
            finish
        fi
    fi
}


function askLogin(){
    echo -e $CYAN""
    echo "###################################################################"
    echo -e "# "$UNDERLINE$YELLOW"WARNING$NORMAL$CYAN: This section of the script CAN destroy your characters #"
    echo "# and accounts information.                                       #"
    echo "# Read questions carefully before you reply.                      #"
    echo "###################################################################"
    echo ""
    echo "Choose full (f) if you don't have and 'accounts' table or would"
    echo "prefer to erase the existing accounts information."
    echo "Choose skip (s) to skip login server DB installation and go to"
    echo "game server DB installation/upgrade."
    echo ""
    echo -ne "Login server DB install type: (f) full, (s) skip or (q) quit? "$ENDC
    read LSP

    case "$LSP" in
    	"f"|"F") loginClean; loginInstall; gameBackup; askGame;;
    	"s"|"S") gameBackup; askGame;;
    	"q"|"Q") finish;;
    	*) askLogin;;
    esac
}


function loginClean(){
    echo -e $YELLOW"Deleting login server tables."$ENDC
    $MY_LOGIN < ls_cleanup.sql
}


function loginInstall(){
    echo ""
    echo -e $CYAN"Installling new loginserver content."$ENDC

    for login in $(ls ./sql/login/*.sql);do
	   echo -e $CYAN"Installing login server table: $login"$ENDC

	   $MY_LOGIN < $login
    done
}


function gameBackup(){
    while true; do
        echo ""
        echo -ne $CYAN"Do you want to make a backup copy of your Game Server DB? (Y/n): "$ENDC
        read GSB

        if [ "$GSB" == "Y" -o "$GSB" == "y" ]; then
            echo -e $CYAN"Making a backup of the original game server database."$ENDC
            mkdir -p $BACKUPS
            DATETIME=$(date '+%Y%m%d_%H%M%S')
            $MYSQLDUMP_PATH --login-path=$GS_DB --add-drop-table $GS_DB > $BACKUPS/GS_Backup_$DATETIME.sql

            if [ $? -ne 0 ]; then
                echo ""
                echo -e $RED "There was a problem accesing your game server database, "
                echo -e "either it wasnt created or authentication data is incorrect."$ENDC
                exit 1
            fi
            echo -e $GREEN"Backup completed."$ENDC
            break
        elif [ "$GSB" == "n" -o "$GSB" == "N" ]; then 
            break
        fi
    done 
}


function loginBackup(){
    while true; do
        echo ""
        echo -ne $CYAN"Do you want to make a backup copy of your Login Server DB? (Y/n): "$ENDC
        read LSB

        if [ "$LSB" == "Y" -o "$LSB" == "y" ]; then
            echo -e $CYAN"Making backup of original login server database."$ENDC
            mkdir -p $BACKUPS
            DATETIME=$(date '+%Y%m%d_%H%M%S')
            $MYSQLDUMP_PATH --login-path=$LS_DB --add-drop-table $LS_DB > $BACKUPS/LS_Backup_$DATETIME.sql

            if [ $? -ne 0 ]; then
                echo ""
                echo -e $RED"There was a problem accesing your login server database, "
                echo -e "either it wasn't created or authentication data is incorrect."$ENDC
                exit 1
            fi
            echo -e $GREEN"Backup completed."$ENDC
            break
        elif [ "$LSB" == "n" -o "$LSB" == "N" ]; then 
            break
        fi
    done 
}


function askGame(){
    echo ""
    echo -e $UNDERLINE$YELLOW"WARNING$NORMAL$CYAN: A full install (f) will destroy all existing character data."
    echo -ne "Game server DB install type: (f) full install, (u) upgrade, (s) skip or (q) quit? "$ENDC
    read GSP

    case "$GSP" in
    	"f"|"F") gameClean; gameInstall I; customTables;;
    	"u"|"U") gameInstall U; customTables;;
    	"s"|"S") customTables;;
    	"q"|"Q") finish;;
    	*) askGame;;
    esac
}

function gameClean(){
    echo -e $YELLOW"Deleting all gameserver tables."$ENDC
    $MY_GAME < gs_cleanup.sql
}

function gameInstall(){
    echo ""
    if [ "$1" == "I" ]; then 
        echo -e $CYAN"Installling new gameserver content."$ENDC
    else
        echo -e $CYAN"Upgrading gameserver content."$ENDC
        echo -e $YELLOW"Not working at now."$ENDC
    fi

    for game in $(ls ./sql/server/*.sql);do
	   echo -e $CYAN"Installing GameServer table: $game"$ENDC

	   $MY_GAME < $game
    done

    # newbie_helper
}

customTables(){
    echo ""
    echo -ne $CYAN"Install custom gameserver DB tables: (y) yes or (n) no or (q) quit? "$ENDC
    #read CUP # Auto skipp this step.
    CUP="Q"
    case "$CUP" in
    	"y"|"Y") cstinstall;;
    	"n"|"N") finish;;
    	"q"|"Q") finish;;
    	*) customTables;;
    esac

    finish
}

finish(){
    echo -e $CYAN
    echo "##################################################################"
    echo "# Script execution finished.                                     #"
    echo "#                                                                #"
    echo "# © 2016-2017 Cute Potato Development.                           #"
    echo -e "# database_installer script comes with "$UNDERLINE"ABSOLUTELY NO WARRANTY"$NORMAL".   #"
    echo "# This software is free, and you are welcome to redistribute it  #"
    echo "# under license conditions.                                      #"
    echo "# See the file license.txt for further details.                  #"
    echo "#                                                                #"
    echo "# Thanks for using our source code.                              #"
    echo "# visit our web for more info about our project.                 #"
    echo "##################################################################"
    echo -e $ENDC
    exit 0
}

# Check or all paths are in place.
function checkMysql(){
    # Check or mysql, mysqldump and mysql_config_editor exists.
    MYSQLDUMP_PATH=$(which -a mysqldump 2> /dev/null)
    MYSQL_PATH=$(which -a mysql 2> /dev/null)
    MYSQLCONFE_PATH=$(which -a mysql_config_editor 2> /dev/null)

    if [ $? -ne 0 ]; then
        echo -e $RED"Unable to find MySQL binaries on path."$ENDC

        while true; do
            echo -ne "\nPlease enter MySQL binaries directory (no trailing slash): "

            read MYSQLBIN_PATH

            if [ -e "$MYSQLBIN_PATH" ] && [ -d "$MYSQLBIN_PATH" ] && [ -e "$MYSQLBIN_PATH/mysqldump" ] && [ -e "$MYSQLBIN_PATH/mysql" ] && [ -e "$MYSQLBIN_PATH/mysql_config_editor" ]; then
                MYSQLDUMP_PATH="$MYSQLBIN_PATH/mysqldump"
                MYSQL_PATH="$MYSQLBIN_PATH/mysql"
                MYSQLCONFE_PATH="$MYSQLBIN_PATH/mysql_config_editor"
                break
            else
               echo "Entered data is invalid. Please verify and try again."
               exit 1
            fi
         done
    fi

    echo -e $CYAN"MySQL paths are ok."$ENDC
    echo ""
}

clear
checkMysql
checkConfig

MY_LOGIN="$MYSQL_PATH --login-path=$LS_DB $LS_DB"
MY_GAME="$MYSQL_PATH --login-path=$GS_DB $GS_DB"

# loginBackup
askLogin
