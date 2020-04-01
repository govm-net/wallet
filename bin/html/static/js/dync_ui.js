function dynamic_load(chain, ui_id, rst) {
    $("#" + ui_id).html("")
    var uis = $("#" + ui_id)
    var nat_list = $('<ul class="nav nav-tabs">')
    console.log("dynamic_load:", rst)
    rst.app = dataEncode(rst.app, "bytes2hex")
    $.each(rst.list, function (i, item) {
        var li = $('<li>').append(
            $('<a href="#uis_id_' + i + '" data-toggle="tab">').append($('<b>').append(item.info.title))
        )
        if (i == 0) {
            li.attr("class", "active")
        }
        // if (item.info.is_view_ui == true) {
        //     li.css('background-color', '#00FFCC')
        // }
        nat_list.append(li)
    });
    uis.append(nat_list)
    var tab_content = $('<div class="tab-content">')
    $.each(rst.list, function (i, item) {
        var it;
        var form = $('<form class="bs-example bs-example-form" role="form">');
        if (i == 0) {
            it = $('<div class="tab-pane active" id="uis_id_' + i + '">')
        } else {
            it = $('<div class="tab-pane" id="uis_id_' + i + '">')
        }
        var desc = $('<div class="input-group">').append(
            $('<span class="input-group-addon" data-localize="description">Description</span>')
        )
        if (item.info.is_view_ui == true) {
            desc.append($('<span class="form-control" data-localize="rui_desc">').append(
                '(Free)View UI,only read and show data.'))
        } else {
            desc.append($('<span class="form-control" data-localize="wui_desc">').append(
                'Run UI,set paramete to run app'))
        }
        form.append(desc)
        if (item.info.is_view_ui != true) {
            if (item.info.cost === undefined) {
                item.info.cost = 0
            }
            form.append($('<div class="input-group">').append(
                $('<span class="input-group-addon" data-localize="energy">Energy</span>')
            ).append(
                $('<input type="number" class="form-control" name="energy" value="1">')
            ).append(
                $('<span class="input-group-addon">t9</span>')
            ));
            form.append($('<div class="input-group">').append(
                $('<span class="input-group-addon" data-localize="cost">Cost</span>')
            ).append(
                $('<input type="text" class="form-control" name="cost" value="' + item.info.cost + '">')
            ));
            if (item.info.value_type == "string") {
                form.append($('<div class="input-group">').append(
                    $('<span class="input-group-addon" data-localize="paramete">Paramete</span>')
                ).append(
                    $('<input type="text" class="form-control" name="paramete">')
                ));
            }
        }

        $.each(item.items, function (j, sub_it) {
            var sit = $('<div class="input-group">').append(
                $('<span class="input-group-addon">').append(sub_it.key)
            )
            var vit = $('<input type="text" class="form-control">')
            if (item.info.is_view_ui == true || sub_it.name === undefined)
                vit.attr("name", "name" + j);
            else
                vit.attr("name", sub_it.name);
            if (sub_it.value !== undefined) {
                vit.attr("value", sub_it.value)
            }
            if (sub_it.value_type !== undefined) {
                switch (sub_it.value_type) {
                    case "number":
                        vit.attr("type", "number")
                    case "bool":
                        vit.attr("type", "checkbox")
                    case "user":
                        vit.attr("value", rst.user)
                }
            }
            if (sub_it.description !== undefined) {
                vit.attr("placeholder", sub_it.description)
                vit.attr("title", sub_it.description)
            }
            if (sub_it.readonly == true) {
                vit.attr("readonly", "readonly")
            }
            if (sub_it.hide == true) {
                sit.attr("style", "display: none")
            }
            sit.append(vit)
            form.append(sit)
        });
        form.append('<br>')
        var btns = $('<div class="pull-right">')
        var next = $('<button type="button" class="btn btn-success">').append("Next")
        if (item.info.is_view_ui == true) {
            btns.append(next)
        }
        var btn = $('<button type="button" class="btn btn-success">')
        if (item.info.is_view_ui == true) {
            btn.html("Search")
        } else {
            btn.html("Run")
        }
        btns.append(btn)
        form.append(btns)
        form.append($('<br>'))
        var rstEle = $('<div>')
        var preKey = ""
        form.append(rstEle)
        loadLanguage()
        btn.on('click', function () {
            console.log("click item:", i)
            if (item.info.is_view_ui == true) {
                var data = $(this).parent().parent("form").serializeJSON();
                var key = "";
                $.each(item.items, function (j, sub_it) {
                    var type = sub_it.value_type;
                    key += dataEncode(data["name" + j], type)
                });
                rstEle.html($('<h4>').append("Read Result,key:" + key));
                preKey = key;
                readFromDB(chain, rst.app, item.info.struct_name, (!item.info.is_log).toString(),
                    key, item.info.value_type, rstEle, item.view_items)
            } else {
                var data = $(this).parent().parent("form").serializeJSON()
                var cost = dataEncode(data.cost, "cost2num");
                var energy = parseInt(1000000000 * data.energy);
                var prifix = "";
                delete data.energy;
                delete data.cost;

                if (item.info.prefix !== undefined) {
                    prifix = item.info.prefix
                }
                if (item.info.value_type == "bytes") {
                    var key = "";
                    $.each(item.items, function (j, sub_it) {
                        key += dataEncode(data["name" + j], sub_it.value_type);
                    });
                    prifix += key;
                    data = null;
                } else if (item.info.value_type == "string") {
                    prifix += dataEncode(data["paramete"], "str2hex");
                    data = null;
                }

                rstEle.html($('<h4>').append("Run Result:"));
                runApp(chain, rst.app, cost, prifix, item.info.value_type, data, rstEle, energy)
            }
        })
        next.on('click', function () {
            getNextKey(chain, rst.app, item.info.struct_name, (!item.info.is_log).toString(),
                preKey, function (next_key) {
                    if (next_key == "") {
                        rstEle.html($('<h4>').append("Read Result:"));
                        rstEle.append("not found next key,Previous Key:" + preKey)
                        preKey = "";
                        return
                    }
                    rstEle.html($('<h4>').append("Read Result,Key:" + next_key));
                    preKey = next_key;
                    readFromDB(chain, rst.app, item.info.struct_name, (!item.info.is_log).toString(),
                        next_key, item.info.value_type, rstEle, item.view_items)
                })
        })
        tab_content.append(it.append(form))
    });
    uis.append(tab_content)
    // $("#" + ui_id).append(uis)
}

function runApp(chain, app, cost, prefix, typ, data, element, energy) {
    var body = { "cost": cost, "energy": energy, "app_name": app, "param": prefix, "param_type": typ, "json_param": data };
    $.ajax({
        type: "POST",
        url: "/api/v1/" + chain + "/transaction/app/run",
        contentType: "application/json",
        data: JSON.stringify(body),
        dataType: "json",
        success: function (rst) {
            console.log(rst);
            element.append(
                $('<span class="label label-success">').append(
                    $('<a href="transaction.html?key=' + rst.trans_key + '&chain=' + chain + '">').append("Transaction:" + rst.trans_key)
                ))
        },
        error: function (XMLHttpRequest, textStatus, errorThrown) {
            console.log(XMLHttpRequest.responseText);
            console.log(textStatus);
            element.append("<span class=\"labellabel-danger\">errorcode:" + XMLHttpRequest.status +
                ". msg: " + XMLHttpRequest.responseText + "</span>");
        }
    });
}

function getNextKey(chain, app, struct, is_db, pre_key, cb) {
    $.ajax({
        type: "GET",
        url: "/api/v1/" + chain + "/data/visit",
        data: { "app_name": app, "struct_name": struct, "is_db_data": is_db, "pre_key": pre_key },
        dataType: "json",
        success: function (rst) {
            if (rst.key === undefined || rst.key == "") {
                cb("")
                return
            }
            cb(rst.key)
        },
        error: function (XMLHttpRequest, textStatus, errorThrown) {
            console.log(XMLHttpRequest.responseText);
            console.log(textStatus);
            cb("")
        }
    });
}

function readFromDB(chain, app, struct, is_db, key, valueType, element, views) {
    $.ajax({
        type: "GET",
        url: "/api/v1/" + chain + "/data",
        data: { "app_name": app, "struct_name": struct, "is_db_data": is_db, "key": key },
        dataType: "json",
        success: function (rst) {
            if (rst.value === undefined || rst.value == "") {
                element.append("not found")
                return
            }
            rst.value = dataEncode(rst.value, valueType)
            rst.life = dataEncode(rst.life, "time2str")
            element.append(
                $('<div class="input-group">').append(
                    $('<span class="input-group-addon">').append("Data Life:")
                ).append(
                    $('<span class="form-control">').append(rst.life)))
            $.each(views, function (j, sub_it) {
                var sit = $('<div class="input-group">').append(
                    $('<span class="input-group-addon">').append(sub_it.key)
                )
                var vit = $('<span class="form-control">').append(
                    dataToString(rst.value, sub_it.vk, sub_it.value_type))
                element.append(sit.append(vit))
            });
        },
        error: function (XMLHttpRequest, textStatus, errorThrown) {
            console.log(XMLHttpRequest.responseText);
            console.log(textStatus);
            element.append("<span class=\"label label-danger\">error code:" + XMLHttpRequest.status +
                ". msg: " + XMLHttpRequest.responseText + "</span>");
        }
    });
}

function dataToString(input, offset, typ) {
    offset = offset || "";
    sp = offset.split(".")
    for (i in sp) {
        if (sp[i] == "") {
            continue
        }
        input = input[sp[i]]
    }
    return dataEncode(input, typ)
}
